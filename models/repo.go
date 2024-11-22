package models

import "time"

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