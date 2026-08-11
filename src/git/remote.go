package git

import (
	"fmt"

	gitlib "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
)

// GetRemoteURL returns the first configured URL for the given remote.
func GetRemoteURL(folder, remoteName string) (string, error) {
	repo, err := gitlib.PlainOpen(folder)
	if err != nil {
		return "", err
	}

	remote, err := repo.Remote(remoteName)
	if err != nil {
		return "", err
	}

	urls := remote.Config().URLs
	if len(urls) == 0 {
		return "", fmt.Errorf("remote '%s' has no URL configured", remoteName)
	}

	return urls[0], nil
}

// HasRemote reports whether a remote with the given name is configured.
func HasRemote(folder, remoteName string) bool {
	url, err := GetRemoteURL(folder, remoteName)
	return err == nil && url != ""
}

// SetRemoteURL creates or replaces a remote pointing at url.
func SetRemoteURL(folder, remoteName, url string) error {
	repo, err := gitlib.PlainOpen(folder)
	if err != nil {
		return err
	}

	// If the remote already exists, drop it first so we can (re)create it
	// with the new URL. CreateRemote fails outright if one already exists.
	_ = repo.DeleteRemote(remoteName)

	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: remoteName,
		URLs: []string{url},
	})

	return err
}
