package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

// getConfigFilePath is a variable to allow overriding in tests
var getConfigFilePath = func() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configPath := filepath.Join(homePath, configFileName)
	return configPath, nil
}

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

// Config method that sets the username passed as the current username in the state's config and
// writes it to the config file
func (cfg *Config) SetUser(name string) error {
	cfg.CurrentUserName = name
	return write(*cfg)
}

// Function that reads the config file and extracts the json data as a Config struct
func Read() (Config, error) {
	fullPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, fmt.Errorf("failed to get path to home: %w", err)
	}

	data, err := os.Open(fullPath)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config file: %w", err)
	}
	defer data.Close()

	newConfig := Config{}
	err = json.NewDecoder(data).Decode(&newConfig)
	if err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config data: %w", err)
	}

	return newConfig, nil
}

// Function that encodes the config struct into json and writes it in the config file
func write(cfg Config) error {
	fullPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}

	defer file.Close()

	err = json.NewEncoder(file).Encode(cfg)
	if err != nil {
		return err
	}

	return nil
}
