package git

import (
	"gitup/src/helper"

	gitlib "github.com/go-git/go-git/v5"
)

func Commit(folder string, files []string, message string) error {
	repo, err := gitlib.PlainOpen(folder)
	if err != nil {
		return err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	status, err := worktree.Status()
	if err != nil {
		return err
	}

	if len(status) == 0 {
		helper.Log("There are no changes to commit.")
		return nil
	}

	helper.Log("Changes to be committed:")

	for file, fileStatus := range status {
		helper.Log("  %c %s", fileStatus.Worktree, file)
	}

	helper.Log("Commit message: %s", message)

	response := helper.ConfirmationDiag(
		"GitUp is about to create this commit.",
		"These changes will be saved to local Git history.",
		"Continue?",
	)

	if response == helper.N {
		helper.Log("Commit cancelled.")
		return nil
	}

	if len(files) == 0 {
		helper.Log("Staging all changes.")

		err = worktree.AddWithOptions(&gitlib.AddOptions{
			All: true,
		})
		if err != nil {
			return err
		}
	} else {
		for _, file := range files {
			helper.Log("Staging %s.", file)

			_, err = worktree.Add(file)
			if err != nil {
				return err
			}
		}
	}

	_, err = worktree.Commit(message, &gitlib.CommitOptions{})
	if err != nil {
		return err
	}

	helper.LogOk("Commit created: %s", message)

	return nil
}
