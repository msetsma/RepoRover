package core

import (
	"errors"
	"fmt"
	"os"
	"io/fs"
	"path/filepath"

	"github.com/pelletier/go-toml"
)

type Config struct {
	DefaultGroup string `toml:"default_group"`
}

const repoRoverDir = ".reporover"

// basePath dynamically resolves to the user's home directory
var basePath string
var defaultConfigFile string

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("Failed to resolve user home directory")
	}

	basePath = filepath.Join(homeDir, repoRoverDir, "groups")
	defaultConfigFile = filepath.Join(homeDir, repoRoverDir, "config.toml")
}

// CreateGroup initializes a new group with a SQLite database and a TOML config file.
func CreateGroup(groupName string) error {
	groupPath := filepath.Join(basePath, groupName)
	dbPath := filepath.Join(groupPath, groupName+".db")
	configPath := filepath.Join(groupPath, groupName+".toml")

	// Ensure base directory exists
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		if err := os.MkdirAll(basePath, 0755); err != nil {
			return fmt.Errorf("failed to create base path: %w", err)
		}
	}

	// Create group directory
	if _, err := os.Stat(groupPath); !os.IsNotExist(err) {
		return errors.New("group already exists")
	}
	if err := os.Mkdir(groupPath, 0755); err != nil {
		return fmt.Errorf("failed to create group path: %w", err)
	}

	// Create SQLite database
	_, err := os.Create(dbPath)
	if err != nil {
		return fmt.Errorf("failed to create SQLite database: %w", err)
	}

	// Create TOML config file
	_, err = os.Create(configPath)
	if err != nil {
		return fmt.Errorf("failed to create group config file: %w", err)
	}

	return nil
}

// SetDefaultGroup updates the global configuration to set the default group.
func SetDefaultGroup(groupName string) error {
	config := &Config{}

	// Read existing config
	if _, err := os.Stat(defaultConfigFile); err == nil {
		data, err := os.ReadFile(defaultConfigFile)
		if err != nil {
			return fmt.Errorf("failed to read config file: %w", err)
		}
		err = toml.Unmarshal(data, config)
		if err != nil {
			return fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	// Update default group
	config.DefaultGroup = groupName
	data, err := toml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to serialize config: %w", err)
	}

	// Write back to file
	err = os.WriteFile(defaultConfigFile, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// ListGroups scans the base path for directories and returns a list of group names.
func ListGroups() ([]string, error) {
	var groups []string

	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		return groups, nil
	}

	// Walk the base directory and collect group names
	err := filepath.WalkDir(basePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Only include directories at the base level
		if d.IsDir() && filepath.Dir(path) == basePath {
			groups = append(groups, d.Name())
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list groups: %w", err)
	}

	return groups, nil
}

// DeleteGroup removes the group directory and all associated files.
func DeleteGroup(groupName string) error {
	groupPath := filepath.Join(basePath, groupName)

	// Check if the group exists
	if _, err := os.Stat(groupPath); os.IsNotExist(err) {
		return fmt.Errorf("group '%s' does not exist", groupName)
	}

	// Remove the group directory and its contents
	err := os.RemoveAll(groupPath)
	if err != nil {
		return fmt.Errorf("failed to delete group '%s': %w", groupName, err)
	}

	// Check if the group is the default group and unset it
	config := &Config{}
	if _, err := os.Stat(defaultConfigFile); err == nil {
		data, err := os.ReadFile(defaultConfigFile)
		if err == nil {
			_ = toml.Unmarshal(data, config)
			if config.DefaultGroup == groupName {
				config.DefaultGroup = ""
				newData, _ := toml.Marshal(config)
				_ = os.WriteFile(defaultConfigFile, newData, 0644)
			}
		}
	}

	return nil
}

// GetDefaultGroup fetches the current default group from the global configuration file.
func GetDefaultGroup() (string, error) {
	config := &Config{}

	// Check if the config file exists
	if _, err := os.Stat(defaultConfigFile); os.IsNotExist(err) {
		return "", nil // No config file, no default group set
	}

	// Read and parse the config file
	data, err := os.ReadFile(defaultConfigFile)
	if err != nil {
		return "", fmt.Errorf("failed to read config file: %w", err)
	}

	err = toml.Unmarshal(data, config)
	if err != nil {
		return "", fmt.Errorf("failed to parse config file: %w", err)
	}

	return config.DefaultGroup, nil
}

// PrintConfig outputs the entire config.toml file to stdout
func PrintConfig() error {
	// Check if the config file exists
	if _, err := os.Stat(defaultConfigFile); os.IsNotExist(err) {
		return fmt.Errorf("config file '%s' does not exist", defaultConfigFile)
	}

	// Read and print the file contents
	data, err := os.ReadFile(defaultConfigFile)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	fmt.Println(string(data))
	return nil
}