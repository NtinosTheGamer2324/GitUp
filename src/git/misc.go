package git

import (
	"os"
	"path/filepath"
)

func IsDirGitProject(folder string) bool {
	_, err := os.Stat(filepath.Join(folder, ".git"))

	return err == nil
}

// ProjectName returns the folder name of a project, used as the default
// name when creating a GitHub repository (e.g. "gitup repo create").
func ProjectName(folder string) (string, error) {
	abs, err := filepath.Abs(folder)
	if err != nil {
		return "", err
	}

	return filepath.Base(abs), nil
}
