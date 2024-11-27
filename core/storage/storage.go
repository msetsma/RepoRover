package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// Database struct to manage the SQLite connection
type Database struct {
	db   *sql.DB
	once sync.Once
}

var (
	instances     = make(map[string]*Database) // Map to manage multiple databases
	instancesLock sync.Mutex                   // Protects the instances map
)

// GetDatabaseInstance returns a singleton instance of the database for the given dbName.
func GetDatabaseInstance(dbName string) (*Database, error) {
	instancesLock.Lock()
	defer instancesLock.Unlock()

	// If an instance already exists for the given dbName, return it
	if instance, exists := instances[dbName]; exists {
		return instance, nil
	}

	// Create a new instance
	instance := &Database{}
	dbPath := getDatabasePath(dbName)
	if err := instance.initDB(dbPath); err != nil {
		return nil, err
	}

	// Save the instance in the map
	instances[dbName] = instance
	return instance, nil
}

// getDatabasePath determines the path to the SQLite database for the given dbName
func getDatabasePath(dbName string) string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "./" // Fallback to current directory
	}
	dbDir := filepath.Join(configDir, ".reporover", "db")
	return filepath.Join(dbDir, dbName+".sqlite")
}

// initDB initializes the SQLite database
func (d *Database) initDB(dbPath string) error {
	// Ensure the directory for the database exists
	if err := os.MkdirAll(filepath.Dir(dbPath), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open the database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	// Verify the database connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}

	// Initialize schema
	if err := d.initSchema(db); err != nil {
		return fmt.Errorf("error initializing schema: %w", err)
	}

	d.db = db
	return nil
}

// initSchema ensures required tables are created
func (d *Database) initSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS repositories (
		id TEXT PRIMARY KEY,
		name TEXT,
		default_branch TEXT,
		remote_url TEXT,
		last_updated DATETIME
	);
	CREATE TABLE IF NOT EXISTS commits (
		id TEXT PRIMARY KEY,
		repository_id TEXT,
		date DATETIME NOT NULL,
		FOREIGN KEY(repository_id) REFERENCES repositories(id)
	);
	`
	_, err := db.Exec(schema)
	return err
}

// Close closes a single database connection
func (d *Database) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

// Close closes all database connections
func CloseAllDatabases() error {
	instancesLock.Lock()
	defer instancesLock.Unlock()

	for name, instance := range instances {
		if err := instance.Close(); err != nil {
			return fmt.Errorf("error closing database '%s': %w", name, err)
		}
		delete(instances, name)
	}
	return nil
}

// SaveRepository saves or updates a repository in the database
func (d *Database) SaveRepository(repo *models.Repository) error {
	query := `
	INSERT INTO repositories (id, name, default_branch, remote_url, last_updated)
	VALUES (?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		name=excluded.name,
		default_branch=excluded.default_branch,
		remote_url=excluded.remote_url,
		last_updated=excluded.last_updated
	`
	_, err := d.db.Exec(query, repo.ID, repo.Name, repo.DefaultBranch, repo.RemoteURL, repo.LastUpdated.Format(time.RFC3339))
	return err
}

// GetRepositories retrieves all repositories from the database
func (d *Database) GetRepositories() ([]models.Repository, error) {
	query := `
	SELECT id, name, default_branch, remote_url, last_updated
	FROM repositories
	`
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying repositories: %w", err)
	}
	defer rows.Close()

	var repos []models.Repository
	for rows.Next() {
		var repo models.Repository
		var lastUpdated string
		if err := rows.Scan(&repo.ID, &repo.Name, &repo.DefaultBranch, &repo.RemoteURL, &lastUpdated); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		repo.LastUpdated, err = time.Parse(time.RFC3339, lastUpdated)
		if err != nil {
			return nil, fmt.Errorf("error parsing last updated timestamp: %w", err)
		}
		repos = append(repos, repo)
	}
	return repos, nil
}

// GetStaleRepositories retrieves repositories not updated for 6 months
func (d *Database) GetStaleRepositories() ([]models.Repository, error) {
	query := `
	SELECT id, name, default_branch, remote_url, last_updated
	FROM repositories
	WHERE last_updated < date('now', '-6 months')
	`
	return d.queryRepositories(query)
}

// GetMostActiveRepositories retrieves repositories with the most commits in the last 30 days
func (d *Database) GetMostActiveRepositories() ([]models.ActiveRepository, error) {
	query := `
	SELECT r.id, r.name, COUNT(c.id) as activity_count
	FROM commits c
	JOIN repositories r ON c.repository_id = r.id
	WHERE c.date >= date('now', '-30 days')
	GROUP BY r.id, r.name
	ORDER BY activity_count DESC
	`
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying most active repositories: %w", err)
	}
	defer rows.Close()

	var activeRepos []models.ActiveRepository
	for rows.Next() {
		var repo models.ActiveRepository
		if err := rows.Scan(&repo.ID, &repo.Name, &repo.ActivityCount); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		activeRepos = append(activeRepos, repo)
	}
	return activeRepos, nil
}

// queryRepositories is a helper for repository queries
func (d *Database) queryRepositories(query string) ([]models.Repository, error) {
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying repositories: %w", err)
	}
	defer rows.Close()

	var repos []models.Repository
	for rows.Next() {
		var repo models.Repository
		var lastUpdated string
		if err := rows.Scan(&repo.ID, &repo.Name, &repo.DefaultBranch, &repo.RemoteURL, &lastUpdated); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		repo.LastUpdated, err = time.Parse(time.RFC3339, lastUpdated)
		if err != nil {
			return nil, fmt.Errorf("error parsing last updated timestamp: %w", err)
		}
		repos = append(repos, repo)
	}
	return repos, nil
}
