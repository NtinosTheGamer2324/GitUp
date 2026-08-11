package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const createRepoURL = "https://api.github.com/user/repos"

// Repository is the subset of GitHub's repo object GitUp cares about.
type Repository struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	HTMLURL  string `json:"html_url"`
	CloneURL string `json:"clone_url"`
	SSHURL   string `json:"ssh_url"`
	Private  bool   `json:"private"`
}

type createRepoRequest struct {
	Name        string `json:"name"`
	Private     bool   `json:"private"`
	Description string `json:"description,omitempty"`
}

type apiError struct {
	Message string `json:"message"`
	Errors  []struct {
		Field   string `json:"field"`
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
}

// CreateRepository creates a new repository under the authenticated user's
// GitHub account.
func CreateRepository(accessToken, name string, private bool, description string) (*Repository, error) {
	body, err := json.Marshal(createRepoRequest{
		Name:        name,
		Private:     private,
		Description: description,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", createRepoURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		var repository Repository
		if err := json.NewDecoder(resp.Body).Decode(&repository); err != nil {
			return nil, err
		}
		return &repository, nil
	}

	var apiErr apiError
	_ = json.NewDecoder(resp.Body).Decode(&apiErr)

	if len(apiErr.Errors) > 0 {
		return nil, fmt.Errorf("%s", apiErr.Errors[0].Message)
	}
	if apiErr.Message != "" {
		return nil, fmt.Errorf("%s", apiErr.Message)
	}

	return nil, fmt.Errorf("GitHub API returned HTTP %d", resp.StatusCode)
}
