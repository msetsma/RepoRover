package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/msetsma/RepoRover/models"
)

// FetchAdditionalRepoData fetches detailed information about a repository
func FetchAdditionalRepoData(org, project, pat, repoID string) (*models.Repository, error) {
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

	var repo models.Repository
	if err := json.NewDecoder(resp.Body).Decode(&repo); err != nil {
		return nil, err
	}

	// Set last updated timestamp
	repo.LastUpdated = time.Now()

	return &repo, nil
}

// FetchRepositories retrieves repositories from Azure DevOps
func FetchRepositories(org, project, pat string) ([]models.Repository, error) {
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
	var reposResponse models.RepositoriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&reposResponse); err != nil {
		return nil, err
	}

	return reposResponse.Value, nil
}