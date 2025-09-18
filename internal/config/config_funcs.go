package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

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

func (cfg *Config) SetUser(name string) error {
	cfg.CurrentUserName = name
	return write(*cfg)
}

func getConfigFilePath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configPath := filepath.Join(homePath, configFileName)
	return configPath, nil
}

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
