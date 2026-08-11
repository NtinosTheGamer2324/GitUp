package git

import (
	"gitup/src/helper"

	gitlib "github.com/go-git/go-git/v5"
)

func InitFolderWithGit(folder string) error {
	_, err := gitlib.PlainInit(folder, false)

	if err != nil {
		if err == gitlib.ErrRepositoryAlreadyExists {
			helper.LogFail("GitUp: Repository already exists in:%s!\n", folder)
			return err
		}

		return err
	}

	helper.LogOk("GitUp: Successfuly initialized folder %s with git.\n", folder)

	return nil
}
