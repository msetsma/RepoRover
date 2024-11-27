package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Manifest struct {
	ActiveGroup   string            `mapstructure:"active_group"`
	DefaultBranch string            `mapstructure:"default_branch"`
	Concurrency   int               `mapstructure:"concurrency"`
	Aliases       map[string]string `mapstructure:"aliases"`
	Paths         Paths             `mapstructure:"paths"`
	Credentials   Credentials       `mapstructure:"credentials"`
	Integrations  Integrations      `mapstructure:"integrations"`
}

type Paths struct {
	Groups string `mapstructure:"clone_destination"`
	Temp   string `mapstructure:"temp"`
}

type Credentials struct {
	Helper  string `mapstructure:"helper"`
	Timeout int    `mapstructure:"timeout"`
}

type Integrations struct {
	Azure Azure `mapstructure:"azure"`
}

type Azure struct {
	Enabled  bool   `mapstructure:"enabled"`
	URL      string `mapstructure:"url"`
	APIToken string `mapstructure:"api_token"`
}

const ConfigFileName = "rover.yaml"

var DefaultBasePath = GetDefaultConfigDirectory()
var DefaultConfigPath = filepath.Join(DefaultBasePath, ConfigFileName)
var DefaultGroupsDirectory = filepath.Join(DefaultBasePath, "groups")
var DefaultTempDirectory = filepath.Join(DefaultBasePath, "tmp")

func setDefaults(v *viper.Viper) {
	v.SetDefault("active_group", "default")
	v.SetDefault("default_branch", "main")
	v.SetDefault("concurrency", 10)
	v.SetDefault("paths.groups", DefaultGroupsDirectory)
	v.SetDefault("paths.temp", DefaultTempDirectory)
	v.SetDefault("credentials.helper", "cache")
	v.SetDefault("credentials.timeout", 3600)
	v.SetDefault("integrations.azure.enabled", false)
	v.SetDefault("integrations.azure.url", "")
	v.SetDefault("integrations.azure.api_token", "")
}

func GetConfigLocation() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	locations := []string{}

	// Only include XDG_CONFIG_HOME if it is set
	xdgConfig := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfig != "" {
		locations = append(locations, filepath.Join(xdgConfig, "rover", ConfigFileName))
	}

	// Add fallback locations
	locations = append(locations,
		filepath.Join(homeDir, ".config", ConfigFileName),
		filepath.Join(homeDir, ".config", "rover", ConfigFileName),
	)

	for _, path := range locations {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}

// GetDefaultConfigPath returns the default dir for the configuration files.
func GetDefaultConfigDirectory() string {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		homeDir = os.Getenv("USERPROFILE")
	}

	xdgConfig := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfig == "" {
		xdgConfig = filepath.Join(homeDir, ".config")
	}

	xdgConfigPath := filepath.Join(xdgConfig, "rover")
	configPath := filepath.Join(homeDir, ".config", "rover")

	if _, err := os.Stat(xdgConfigPath); err == nil {
		return xdgConfigPath
	}
	return configPath
}

func CreateConfigFile(v *viper.Viper) error {
	// Create the configuration directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(DefaultConfigPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	// Write default values to the config file
	if err := v.WriteConfigAs(DefaultConfigPath); err != nil {
		return fmt.Errorf("failed to write default config file: %w", err)
	}
	return nil
}

func expandEnvVariables(cfg interface{}) interface{} {
	switch v := cfg.(type) {
	case map[string]interface{}:
		for key, value := range v {
			v[key] = expandEnvVariables(value)
		}
		return v
	case []interface{}:
		for i, value := range v {
			v[i] = expandEnvVariables(value)
		}
		return v
	case string:
		return os.ExpandEnv(v)
	default:
		return v
	}
}

// Update saves the given Manifest struct to the configuration file.
func Update(manifest *Manifest) error {
	configPath := GetConfigLocation()

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configMap := make(map[string]interface{})
	err := mapstructure.Decode(manifest, &configMap)
	if err != nil {
		return fmt.Errorf("invalid config structure")
	}

	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	// Load the existing configuration if it exists
	if _, err := os.Stat(configPath); err == nil {
		if err := v.ReadInConfig(); err != nil {
			return fmt.Errorf("failed to read existing config file")
		}
	}

	for key, value := range configMap {
		v.Set(key, value)
	}

	if err := v.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to write updated config file")
	}
	return nil
}

func Load() (*Manifest, error) {
	v := viper.New()
	setDefaults(v)
	configPath := GetConfigLocation()
	v.SetConfigFile(configPath)
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		fmt.Println(err)
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			CreateConfigFile(v)
		}
	}

	var rawConfig map[string]interface{}
	if err := v.Unmarshal(&rawConfig); err != nil {
		return nil, fmt.Errorf("error unmarshalling config to map: %v", err)
	}

	expandedConfig := expandEnvVariables(rawConfig)
	cfg := &Manifest{}
	decoderConfig := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   cfg,
		TagName:  "mapstructure",
	}
	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating decoder: %v", err)
	}
	if err := decoder.Decode(expandedConfig); err != nil {
		return nil, fmt.Errorf("error decoding config into struct: %v", err)
	}

	return cfg, nil
}
