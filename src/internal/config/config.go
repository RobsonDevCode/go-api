package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const FilePath = "C:/Users/RobsonBasquill-Lipsc/Repos/Test/go-api/src/config/config.yaml"

type Config struct {
	Database          DBConfig
	UserClientOptions UserClient
	JWTSettings       JWTSettings
}

type DBConfig struct {
	Driver            string `yaml:"driver"`
	Host              string `yaml:"host"`
	Port              int    `yaml:"port"`
	Name              string `yaml:"name"`
	Username          string `yaml:"username"`
	Password          string `yaml:"password"`
	SSlMode           string `yaml:"sslMode"`
	AllowZeroDateTime bool   `yaml:"allowZeroDateTime"`
}

type UserClient struct {
	URL string `yaml:"baseUrl"`
}

type JWTSettings struct {
	Issuer                string `yaml:"issuer"`
	Audience              string `yaml:"audience"`
	ExpiresInMinutes      int    `yaml:"expiresInMinutes"`
	RefreshTokenExpiresIn int    `yaml:"refreshTokenExpiresIn"`
	Key                   string `yaml:"key"`
}

func Load() (*Config, error) {

	data, err := os.ReadFile(FilePath)
	if err != nil {
		return nil, err
	}
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configL %w", err)
	}

	return &config, nil
}

func validateConfig(config *Config) error {
	if config.Database.Username == "" {
		return fmt.Errorf("database username is required")
	}
	if config.Database.Name == "" {
		return fmt.Errorf("database name is required")
	}

	if config.Database.Port <= 0 {
		return fmt.Errorf("server port must be positive")
	}

	return nil
}
