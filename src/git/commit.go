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
