package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestConfig_SetUser(t *testing.T) {
	tests := []struct {
		name     string
		username string
		config   Config
	}{
		{
			name:     "set user on empty config",
			username: "testuser",
			config:   Config{},
		},
		{
			name:     "set user on existing config",
			username: "newuser",
			config: Config{
				DbUrl:           "postgres://test",
				CurrentUserName: "olduser",
			},
		},
		{
			name:     "set empty username",
			username: "",
			config:   Config{DbUrl: "postgres://test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary config file for testing
			tempDir := t.TempDir()
			tempFile := filepath.Join(tempDir, configFileName)
			
			// Create initial config file
			initialData, _ := json.Marshal(tt.config)
			err := os.WriteFile(tempFile, initialData, 0644)
			if err != nil {
				t.Fatalf("failed to create temp config file: %v", err)
			}

			// Temporarily override the config file path
			originalGetConfigFilePath := getConfigFilePath
			getConfigFilePath = func() (string, error) {
				return tempFile, nil
			}
			defer func() {
				getConfigFilePath = originalGetConfigFilePath
			}()

			// Test SetUser
			config := tt.config
			err = config.SetUser(tt.username)
			if err != nil {
				t.Errorf("SetUser() error = %v", err)
				return
			}

			// Verify the user was set in memory
			if config.CurrentUserName != tt.username {
				t.Errorf("expected CurrentUserName = %q, got %q", tt.username, config.CurrentUserName)
			}

			// Verify the config was written to file
			data, err := os.ReadFile(tempFile)
			if err != nil {
				t.Errorf("failed to read config file: %v", err)
				return
			}

			var savedConfig Config
			err = json.Unmarshal(data, &savedConfig)
			if err != nil {
				t.Errorf("failed to unmarshal saved config: %v", err)
				return
			}

			if savedConfig.CurrentUserName != tt.username {
				t.Errorf("expected saved CurrentUserName = %q, got %q", tt.username, savedConfig.CurrentUserName)
			}
		})
	}
}

func TestRead(t *testing.T) {
	tests := []struct {
		name           string
		setupFile      func(string) error
		expectedConfig Config
		expectError    bool
	}{
		{
			name: "read valid config",
			setupFile: func(path string) error {
				config := Config{
					DbUrl:           "postgres://localhost/test",
					CurrentUserName: "testuser",
				}
				data, _ := json.Marshal(config)
				return os.WriteFile(path, data, 0644)
			},
			expectedConfig: Config{
				DbUrl:           "postgres://localhost/test",
				CurrentUserName: "testuser",
			},
			expectError: false,
		},
		{
			name: "read config with missing fields",
			setupFile: func(path string) error {
				config := Config{
					DbUrl: "postgres://localhost/test",
					// CurrentUserName is missing
				}
				data, _ := json.Marshal(config)
				return os.WriteFile(path, data, 0644)
			},
			expectedConfig: Config{
				DbUrl:           "postgres://localhost/test",
				CurrentUserName: "",
			},
			expectError: false,
		},
		{
			name: "read non-existent config",
			setupFile: func(path string) error {
				// Don't create the file
				return nil
			},
			expectedConfig: Config{},
			expectError:    true,
		},
		{
			name: "read invalid JSON config",
			setupFile: func(path string) error {
				return os.WriteFile(path, []byte("invalid json"), 0644)
			},
			expectedConfig: Config{},
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary config file for testing
			tempDir := t.TempDir()
			tempFile := filepath.Join(tempDir, configFileName)

			// Setup the test file
			err := tt.setupFile(tempFile)
			if err != nil {
				t.Fatalf("failed to setup test file: %v", err)
			}

			// Temporarily override the config file path
			originalGetConfigFilePath := getConfigFilePath
			getConfigFilePath = func() (string, error) {
				return tempFile, nil
			}
			defer func() {
				getConfigFilePath = originalGetConfigFilePath
			}()

			// Test Read
			config, err := Read()

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Read() error = %v", err)
				return
			}

			if config.DbUrl != tt.expectedConfig.DbUrl {
				t.Errorf("expected DbUrl = %q, got %q", tt.expectedConfig.DbUrl, config.DbUrl)
			}

			if config.CurrentUserName != tt.expectedConfig.CurrentUserName {
				t.Errorf("expected CurrentUserName = %q, got %q", tt.expectedConfig.CurrentUserName, config.CurrentUserName)
			}
		})
	}
}

func TestWrite(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{
			name: "write complete config",
			config: Config{
				DbUrl:           "postgres://localhost/test",
				CurrentUserName: "testuser",
			},
		},
		{
			name: "write empty config",
			config: Config{},
		},
		{
			name: "write config with special characters",
			config: Config{
				DbUrl:           "postgres://user:pass@localhost/test?sslmode=disable",
				CurrentUserName: "user-with-dash",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary config file for testing
			tempDir := t.TempDir()
			tempFile := filepath.Join(tempDir, configFileName)

			// Temporarily override the config file path
			originalGetConfigFilePath := getConfigFilePath
			getConfigFilePath = func() (string, error) {
				return tempFile, nil
			}
			defer func() {
				getConfigFilePath = originalGetConfigFilePath
			}()

			// Test write
			err := write(tt.config)
			if err != nil {
				t.Errorf("write() error = %v", err)
				return
			}

			// Verify the file was created and contains correct data
			data, err := os.ReadFile(tempFile)
			if err != nil {
				t.Errorf("failed to read written file: %v", err)
				return
			}

			var savedConfig Config
			err = json.Unmarshal(data, &savedConfig)
			if err != nil {
				t.Errorf("failed to unmarshal written config: %v", err)
				return
			}

			if savedConfig.DbUrl != tt.config.DbUrl {
				t.Errorf("expected DbUrl = %q, got %q", tt.config.DbUrl, savedConfig.DbUrl)
			}

			if savedConfig.CurrentUserName != tt.config.CurrentUserName {
				t.Errorf("expected CurrentUserName = %q, got %q", tt.config.CurrentUserName, savedConfig.CurrentUserName)
			}
		})
	}
}

func TestGetConfigFilePath(t *testing.T) {
	path, err := getConfigFilePath()
	if err != nil {
		t.Errorf("getConfigFilePath() error = %v", err)
		return
	}

	if path == "" {
		t.Errorf("expected non-empty path")
	}

	// Verify the path ends with the config file name
	if filepath.Base(path) != configFileName {
		t.Errorf("expected path to end with %q, got %q", configFileName, filepath.Base(path))
	}

	// Verify the path contains the home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Errorf("failed to get home directory: %v", err)
		return
	}

	expectedPath := filepath.Join(homeDir, configFileName)
	if path != expectedPath {
		t.Errorf("expected path = %q, got %q", expectedPath, path)
	}
}