package git

import (
	"fmt"

	"gitup/src/helper"

	gitlib "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	httptransport "github.com/go-git/go-git/v5/plumbing/transport/http"
)

// Publish pushes the current branch to remoteName (typically "origin"),
// refusing to overwrite remote history GitUp hasn't seen locally.
func Publish(folder, remoteName, username, accessToken string) error {
	repo, err := gitlib.PlainOpen(folder)
	if err != nil {
		return err
	}

	head, err := repo.Head()
	if err != nil {
		return fmt.Errorf("could not determine the current branch: %w", err)
	}

	if !head.Name().IsBranch() {
		return fmt.Errorf("not currently on a branch (detached HEAD)")
	}

	branch := head.Name().Short()

	if _, err := repo.Remote(remoteName); err != nil {
		return fmt.Errorf(
			"no '%s' remote is configured. Run 'gitup repo create' first, or set one manually",
			remoteName,
		)
	}

	auth := &httptransport.BasicAuth{
		Username: username,
		Password: accessToken,
	}

	helper.Log("Checking for remote changes...")

	fetchErr := repo.Fetch(&gitlib.FetchOptions{
		RemoteName: remoteName,
		Auth:       auth,
		Tags:       gitlib.NoTags,
	})

	fetchOk := fetchErr == nil || fetchErr == gitlib.NoErrAlreadyUpToDate
	if !fetchOk {
		// Most likely cause: the remote repo is brand new and empty, so
		// there's nothing to fetch yet. Treat this as a first publish
		// rather than blocking the user.
		helper.Log("Could not check remote history (%s).", fetchErr)
		helper.Log("Assuming this is the first publish to this remote.")
	}

	if fetchOk {
		remoteRefName := plumbing.NewRemoteReferenceName(remoteName, branch)
		remoteRef, refErr := repo.Reference(remoteRefName, true)

		if refErr == nil && remoteRef.Hash() != head.Hash() {
			localCommit, err := repo.CommitObject(head.Hash())
			if err != nil {
				return err
			}

			remoteCommit, err := repo.CommitObject(remoteRef.Hash())
			if err != nil {
				return err
			}

			isAncestor, err := remoteCommit.IsAncestor(localCommit)
			if err != nil {
				return err
			}

			if !isAncestor {
				helper.LogFail("Cannot publish.")
				helper.Log("")
				helper.Log("Origin contains history that does not exist locally.")
				helper.Log("Publishing now would overwrite work GitUp has never seen.")
				helper.Log("")
				helper.Log("Pull the remote changes first, then try again.")
				return nil
			}
		}
	}

	remoteURL, _ := GetRemoteURL(folder, remoteName)

	response := helper.ConfirmationDiag(
		"GitUp is about to publish your changes.",
		fmt.Sprintf("This will push branch '%s' to '%s' (%s).", branch, remoteName, remoteURL),
		"Continue?",
	)

	if response == helper.N {
		helper.Log("Publish cancelled.")
		return nil
	}

	helper.Log("Publishing branch '%s'...", branch)

	refSpec := config.RefSpec(fmt.Sprintf("refs/heads/%s:refs/heads/%s", branch, branch))

	err = repo.Push(&gitlib.PushOptions{
		RemoteName: remoteName,
		Auth:       auth,
		RefSpecs:   []config.RefSpec{refSpec},
	})

	if err == gitlib.NoErrAlreadyUpToDate {
		helper.Log("Already up to date. Nothing to publish.")
		return nil
	}

	if err != nil {
		return fmt.Errorf("failed to publish: %w", err)
	}

	helper.LogOk("Published successfully.")

	return nil
}
