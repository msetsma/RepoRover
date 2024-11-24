package config

import (
	"fmt"
	"os"
	"context"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Config struct {
    ActiveGroup          string               `mapstructure:"active_group"`
    DefaultBranch        string               `mapstructure:"default_branch"`
	Concurrency          int                  `mapstructure:"concurrency"`
    Aliases              map[string]string    `mapstructure:"aliases"`
    Paths                Paths                `mapstructure:"paths"`
    Credentials          Credentials          `mapstructure:"credentials"`
    Integrations         Integrations         `mapstructure:"integrations"`
}

type Paths struct {
    CloneDestination string `mapstructure:"clone_destination"`
    Temp             string `mapstructure:"temp"`
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


func setDefaults(v *viper.Viper) {
    v.SetDefault("active_group", "default")
    v.SetDefault("default_branch", "main")
    v.SetDefault("concurrency", 10)
    v.SetDefault("paths.clone_destination", "~/reporover_projects")
    v.SetDefault("paths.temp", "/tmp/reporover/")
    v.SetDefault("credentials.helper", "cache")
    v.SetDefault("credentials.timeout", 3600)
    v.SetDefault("integrations.azure.enabled", false)
    v.SetDefault("integrations.azure.url", "")
    v.SetDefault("integrations.azure.api_token", "")
}


func LoadConfig() (*Config, error) {
    v := viper.New()

    // Set the config file name (without extension)
    v.SetConfigName("rover") 
    v.SetConfigType("yaml") 
    v.AddConfigPath("$HOME/.rover")
	setDefaults(v)
    v.AutomaticEnv()

    if err := v.ReadInConfig(); err != nil {
        fmt.Errorf("Issue with config file; rover is using default values")
    }

    var rawConfig map[string]interface{}
    if err := v.Unmarshal(&rawConfig); err != nil {
        return nil, fmt.Errorf("error unmarshalling config to map: %v", err)
    }

    expandedConfig := expandEnvVariables(rawConfig)
    cfg := &Config{}
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

// Context helper functions
type key int

const configKey key = iota

func WithConfig(ctx context.Context, cfg *Config) context.Context {
    return context.WithValue(ctx, configKey, cfg)
}

func FromContext(ctx context.Context) *Config {
    cfg, ok := ctx.Value(configKey).(*Config)
    if !ok {
        return nil
    }
    return cfg
}
