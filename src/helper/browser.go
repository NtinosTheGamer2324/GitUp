package helper

import (
	"os/exec"
	"runtime"
)

// OpenURL attempts to open the given URL in the user's default browser.
// It fails silently on purpose: callers should always also print the URL
// so the user can open it manually if this doesn't work.
func OpenURL(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}

	return cmd.Start()
}
