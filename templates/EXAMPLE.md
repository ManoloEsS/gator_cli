# Example: How to Use the Config Test Template

This example shows how to transform the `config_test_template.go` into a working test file.

## Original Template Section:
```go
func TestConfig_SetUser(t *testing.T) {
	// TODO: Define test cases using table-driven test pattern
	tests := []struct {
		name     string // Descriptive name for the test case
		username string // Input username to set
		config   Config // Initial config state
		// TODO: Add more fields as needed (expectedError bool, etc.)
	}{
		{
			// TODO: Add test case for setting user on empty config
			name:     "TODO: add descriptive name",
			username: "TODO: add test username",
			config:   Config{}, // TODO: set initial config state
		},
		// TODO: Add more test cases
	}
```

## Completed Implementation:
```go
func TestConfig_SetUser(t *testing.T) {
	// Completed: Define test cases using table-driven test pattern
	tests := []struct {
		name     string // Descriptive name for the test case
		username string // Input username to set
		config   Config // Initial config state
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
			// Create temporary directory and file for testing
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

			// Call the method being tested
			config := tt.config
			err = config.SetUser(tt.username)
			if err != nil {
				t.Errorf("SetUser() error = %v", err)
				return
			}

			// Verify the result in memory
			if config.CurrentUserName != tt.username {
				t.Errorf("expected CurrentUserName = %q, got %q", tt.username, config.CurrentUserName)
			}

			// Verify the result was persisted to file
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
```

## Key Learning Points Demonstrated:

1. **Table-Driven Tests**: Multiple test cases in a structured format
2. **Temporary File Handling**: Using `t.TempDir()` for isolation
3. **Function Override Pattern**: Temporarily replacing global functions for testing
4. **JSON Testing**: Marshall/unmarshall validation
5. **Error Handling**: Proper error checking and reporting
6. **Resource Cleanup**: Using `defer` for cleanup

## How the Template Helped:

1. **Structure**: Provided the complete test function structure
2. **Comments**: Explained the purpose of each section
3. **Patterns**: Showed Go testing best practices
4. **TODO Guidance**: Clear markers for what to implement
5. **Examples**: Concrete examples of test cases and assertions

This approach allows you to learn by:
- Understanding the pattern before implementing
- Following guided steps with clear markers
- Seeing complete working examples
- Learning best practices through comments and structure