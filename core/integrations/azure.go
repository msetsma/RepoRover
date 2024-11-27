package azure

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Repository struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	DefaultBranch string            `json:"defaultBranch"`
	RemoteURL     string            `json:"remoteUrl"`
	LastUpdated   time.Time         `json:"lastUpdated"`
	Languages     map[string]uint64 `json:"languages"`
}

// ActiveRepository represents a repository with its activity count
type ActiveRepository struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	ActivityCount int    `json:"activity_count"`
}

// RepositoriesResponse represents the response from Azure DevOps API
type RepositoriesResponse struct {
	Value []Repository `json:"value"`
}

// FetchAdditionalRepoData fetches detailed information about a repository
func FetchAdditionalRepoData(org, project, pat, repoID string) (*Repository, error) {
	url := fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/git/repositories/%s?api-version=7.1-preview.1", org, project, repoID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth("", pat)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch repository data: %s", resp.Status)
	}

	var repo Repository
	if err := json.NewDecoder(resp.Body).Decode(&repo); err != nil {
		return nil, err
	}

	// Set last updated timestamp
	repo.LastUpdated = time.Now()

	return &repo, nil
}

// FetchRepositories retrieves repositories from Azure DevOps
func FetchRepositories(org, project, pat string) ([]Repository, error) {
	// Azure DevOps REST API URL
	url := fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/git/repositories?api-version=7.1-preview.1", org, project)

	// Create HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth("", pat)

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch repositories: %s", resp.Status)
	}
	var reposResponse RepositoriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&reposResponse); err != nil {
		return nil, err
	}

	return reposResponse.Value, nil
}
