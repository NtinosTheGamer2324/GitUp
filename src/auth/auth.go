package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// GitHubClientID is the Client ID of the GitUp GitHub App.
// Device flow does not require a client secret.
const GitHubClientID = "Iv23ctEafUjuNuS3k4vZ"

const (
	deviceCodeURL  = "https://github.com/login/device/code"
	accessTokenURL = "https://github.com/login/oauth/access_token"
	userURL        = "https://api.github.com/user"
)

// Credentials is what GitUp stores on disk after a successful login.
type Credentials struct {
	Login             string    `json:"login"`
	AccessToken       string    `json:"access_token"`
	RefreshToken      string    `json:"refresh_token,omitempty"`
	ExpiresAt         time.Time `json:"expires_at,omitempty"`
	RefreshExpiresAt  time.Time `json:"refresh_expires_at,omitempty"`
}

// Expired reports whether the access token is known to be expired.
// Credentials with no expiry information are treated as always valid.
func (c *Credentials) Expired() bool {
	if c.ExpiresAt.IsZero() {
		return false
	}
	// Refresh a little early to avoid races against in-flight requests.
	return time.Now().After(c.ExpiresAt.Add(-30 * time.Second))
}

// RefreshExpired reports whether the refresh token itself has expired,
// meaning the user must run `gitup login` again.
func (c *Credentials) RefreshExpired() bool {
	if c.RefreshExpiresAt.IsZero() {
		return false
	}
	return time.Now().After(c.RefreshExpiresAt)
}

// GitHubUser is the subset of GitHub's /user response GitUp cares about.
type GitHubUser struct {
	Login string `json:"login"`
}

// DeviceCode is GitHub's response to starting the device flow.
type DeviceCode struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

type tokenResponse struct {
	AccessToken           string `json:"access_token"`
	RefreshToken          string `json:"refresh_token"`
	ExpiresIn             int    `json:"expires_in"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
	TokenType             string `json:"token_type"`
	Error                 string `json:"error"`
	ErrorDescription      string `json:"error_description"`
}

// ErrAuthorizationPending is returned by Poll while the user hasn't
// finished authorizing in the browser yet. Callers should keep polling.
var ErrAuthorizationPending = fmt.Errorf("authorization_pending")

// ErrSlowDown is returned when GitHub asks us to poll less frequently.
var ErrSlowDown = fmt.Errorf("slow_down")

// RequestDeviceCode starts the OAuth device flow.
func RequestDeviceCode() (*DeviceCode, error) {
	form := url.Values{}
	form.Set("client_id", GitHubClientID)

	req, err := http.NewRequest("POST", deviceCodeURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var dc DeviceCode
	if err := json.NewDecoder(resp.Body).Decode(&dc); err != nil {
		return nil, err
	}

	if dc.DeviceCode == "" {
		return nil, fmt.Errorf("GitHub did not return a device code (HTTP %d)", resp.StatusCode)
	}

	return &dc, nil
}

// PollOnce checks once whether the user has finished authorizing.
// It returns ErrAuthorizationPending or ErrSlowDown if not done yet.
func PollOnce(deviceCode string) (*Credentials, error) {
	form := url.Values{}
	form.Set("client_id", GitHubClientID)
	form.Set("device_code", deviceCode)
	form.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")

	tok, err := requestToken(form)
	if err != nil {
		return nil, err
	}

	return credentialsFromToken(tok)
}

// RefreshAccessToken exchanges a refresh token for a new access token.
func RefreshAccessToken(refreshToken string) (*Credentials, error) {
	form := url.Values{}
	form.Set("client_id", GitHubClientID)
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", refreshToken)

	tok, err := requestToken(form)
	if err != nil {
		return nil, err
	}

	return credentialsFromToken(tok)
}

func requestToken(form url.Values) (*tokenResponse, error) {
	req, err := http.NewRequest("POST", accessTokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tok tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tok); err != nil {
		return nil, err
	}

	switch tok.Error {
	case "":
		// success
	case "authorization_pending":
		return nil, ErrAuthorizationPending
	case "slow_down":
		return nil, ErrSlowDown
	default:
		msg := tok.ErrorDescription
		if msg == "" {
			msg = tok.Error
		}
		return nil, fmt.Errorf("%s", msg)
	}

	if tok.AccessToken == "" {
		return nil, fmt.Errorf("GitHub did not return an access token")
	}

	return &tok, nil
}

func credentialsFromToken(tok *tokenResponse) (*Credentials, error) {
	creds := &Credentials{
		AccessToken:  tok.AccessToken,
		RefreshToken: tok.RefreshToken,
	}

	if tok.ExpiresIn > 0 {
		creds.ExpiresAt = time.Now().Add(time.Duration(tok.ExpiresIn) * time.Second)
	}
	if tok.RefreshTokenExpiresIn > 0 {
		creds.RefreshExpiresAt = time.Now().Add(time.Duration(tok.RefreshTokenExpiresIn) * time.Second)
	}

	user, err := FetchUser(creds.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("authenticated, but failed to fetch user identity: %w", err)
	}
	creds.Login = user.Login

	return creds, nil
}

// FetchUser looks up the identity of the account behind an access token.
func FetchUser(accessToken string) (*GitHubUser, error) {
	req, err := http.NewRequest("GET", userURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned HTTP %d", resp.StatusCode)
	}

	var user GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

// credentialsPath returns the on-disk location GitUp stores auth state in.
//
// This is plain-file storage, not OS-native secure credential storage
// (Keychain / Credential Manager / Secret Service). That's a known
// follow-up, tracked as a future improvement rather than done here.
func credentialsPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		home, herr := os.UserHomeDir()
		if herr != nil {
			return "", err
		}
		dir = home
	}

	return filepath.Join(dir, "gitup", "credentials.json"), nil
}

// SaveCredentials writes credentials to disk, restricted to the current user.
func SaveCredentials(creds *Credentials) error {
	path, err := credentialsPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

// LoadCredentials reads previously saved credentials, if any.
// It returns (nil, nil) if the user has never logged in.
func LoadCredentials() (*Credentials, error) {
	path, err := credentialsPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, err
	}

	return &creds, nil
}

// DeleteCredentials removes any stored credentials (used by `gitup logout`).
// It is not an error to call this when no credentials exist.
func DeleteCredentials() error {
	path, err := credentialsPath()
	if err != nil {
		return err
	}

	err = os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

// EnsureValidCredentials loads stored credentials and transparently
// refreshes the access token if it's expired but the refresh token isn't.
// It returns (nil, nil) if the user is simply not logged in.
func EnsureValidCredentials() (*Credentials, error) {
	creds, err := LoadCredentials()
	if err != nil {
		return nil, err
	}
	if creds == nil {
		return nil, nil
	}

	if !creds.Expired() {
		return creds, nil
	}

	if creds.RefreshToken == "" || creds.RefreshExpired() {
		return nil, fmt.Errorf("session expired, please run 'gitup login' again")
	}

	refreshed, err := RefreshAccessToken(creds.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh session, please run 'gitup login' again: %w", err)
	}

	if err := SaveCredentials(refreshed); err != nil {
		return nil, err
	}

	return refreshed, nil
}
