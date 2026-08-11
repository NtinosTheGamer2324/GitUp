package auth

import (
	"time"

	"gitup/src/helper"
)

const pollTimeoutSafetyMargin = 10 * time.Second

// Login runs the full GitHub device-flow login and stores the resulting
// credentials on disk.
func Login() error {
	existing, err := LoadCredentials()
	if err == nil && existing != nil && !existing.Expired() {
		helper.Log("Already logged in as %s.", existing.Login)
		helper.Log("Run 'gitup logout' first if you want to switch accounts.")
		return nil
	}

	helper.Log("Starting GitHub login...")

	device, err := RequestDeviceCode()
	if err != nil {
		return err
	}

	helper.LogOk("Almost there.")
	helper.Log("")
	helper.Log("1. Open: %s", device.VerificationURI)
	helper.Log("2. Enter this code: %s", device.UserCode)
	helper.Log("")

	_ = helper.OpenURL(device.VerificationURI)

	interval := device.Interval
	if interval <= 0 {
		interval = 5
	}

	deadline := time.Now().Add(time.Duration(device.ExpiresIn)*time.Second - pollTimeoutSafetyMargin)

	helper.Log("Waiting for you to finish in the browser...")

	for {
		if time.Now().After(deadline) {
			return errExpiredLogin
		}

		time.Sleep(time.Duration(interval) * time.Second)

		creds, err := PollOnce(device.DeviceCode)
		if err == ErrAuthorizationPending {
			continue
		}
		if err == ErrSlowDown {
			interval += 5
			continue
		}
		if err != nil {
			return err
		}

		if err := SaveCredentials(creds); err != nil {
			return err
		}

		helper.LogOk("GitHub account connected.")
		helper.Log("")
		helper.Log("Logged in as: %s", creds.Login)
		helper.Log("Authentication: GitHub OAuth (device flow)")
		helper.Log("")
		helper.Log("You're ready to publish.")

		return nil
	}
}

// errExpiredLogin is returned when the user doesn't finish the device
// flow before GitHub's code expires.
var errExpiredLogin = errExpired{}

type errExpired struct{}

func (errExpired) Error() string {
	return "login timed out before it was completed, please run 'gitup login' again"
}

// Logout removes any stored GitHub credentials.
func Logout() error {
	creds, _ := LoadCredentials()

	if err := DeleteCredentials(); err != nil {
		return err
	}

	if creds != nil {
		helper.LogOk("Logged out of %s.", creds.Login)
	} else {
		helper.Log("You weren't logged in.")
	}

	return nil
}

// WhoAmI prints the currently authenticated GitHub identity.
func WhoAmI() error {
	creds, err := EnsureValidCredentials()
	if err != nil {
		return err
	}

	helper.Log("GitUp Identity")
	helper.Log("")

	if creds == nil {
		helper.Log("GitHub: Not connected")
		helper.Log("")
		helper.Log("Run 'gitup login' to connect a GitHub account.")
		return nil
	}

	helper.Log("GitHub: %s", creds.Login)

	return nil
}
