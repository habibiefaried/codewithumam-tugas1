package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	URL      string `yaml:"db_url"`
	Name     string `yaml:"db_name"`
	User     string `yaml:"db_user"`
	Password string `yaml:"db_password"`
	DBPort   string `yaml:"db_port"`
	Port     string `yaml:"port"`
}

// LoadConfig loads configuration from secrets.yml if it exists, otherwise uses environment variables
func LoadConfig() (*DatabaseConfig, error) {
	cfg := &DatabaseConfig{}

	// Try multiple locations for secrets.yml
	secretsPath := ""

	// 1. First, try current working directory
	if _, err := os.Stat("secrets.yml"); err == nil {
		secretsPath = "secrets.yml"
	} else {
		// 2. Then try the executable directory
		exePath, err := os.Executable()
		if err == nil {
			exeDir := filepath.Dir(exePath)
			potentialPath := filepath.Join(exeDir, "secrets.yml")
			if _, err := os.Stat(potentialPath); err == nil {
				secretsPath = potentialPath
			}
		}
	}

	// Try to load from secrets.yml if found
	if secretsPath != "" {
		data, err := os.ReadFile(secretsPath)
		if err == nil {
			// File exists, parse it
			if err := yaml.Unmarshal(data, cfg); err != nil {
				return nil, fmt.Errorf("failed to parse secrets.yml: %w", err)
			}
		}
	}

	// Fall back to environment variables if values are not set
	if cfg.URL == "" {
		cfg.URL = os.Getenv("DB_URL")
		if cfg.URL == "" {
			cfg.URL = "localhost"
		}
	}

	if cfg.Name == "" {
		cfg.Name = os.Getenv("DB_NAME")
		if cfg.Name == "" {
			cfg.Name = "postgres"
		}
	}

	if cfg.User == "" {
		cfg.User = os.Getenv("DB_USER")
		if cfg.User == "" {
			cfg.User = "postgres"
		}
	}

	if cfg.Password == "" {
		cfg.Password = os.Getenv("DB_PASSWORD")
		if cfg.Password == "" {
			cfg.Password = "postgres"
		}
	}

	if cfg.DBPort == "" {
		cfg.DBPort = os.Getenv("DB_PORT")
		if cfg.DBPort == "" {
			cfg.DBPort = "5432"
		}
	}

	if cfg.Port == "" {
		cfg.Port = os.Getenv("PORT")
		if cfg.Port == "" {
			cfg.Port = "8080"
		}
	}

	return cfg, nil
}
