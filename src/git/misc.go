package git

import (
	"os"
	"path/filepath"
)

func IsDirGitProject(folder string) bool {
	_, err := os.Stat(filepath.Join(folder, ".git"))

	return err == nil
}
