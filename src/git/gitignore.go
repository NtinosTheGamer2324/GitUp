package git

import (
	"fmt"
	"os"

	"gitup/src/helper"
)

var gitignore = "bin/\ndist/\nbuild/\n*.exe\n"

func WriteGitIgnore(folder string) {
	data := []byte(gitignore)

	err := os.WriteFile(fmt.Sprintf("%s/.gitignore", folder), data, 0644)

	if err != nil {
		helper.LogFail("Failed to write .gitignore: %v", err)
		return
	}

	helper.LogOk("Successfully wrote default .gitignore")
}
