# A Detailed Guide to Go Testing with `gator_cli`

This document provides a detailed, line-by-line walkthrough of how to create tests in Go, using the `gator_cli` project as a practical example.

## Part 1: Creating Tests in Go - The Fundamentals

Go's testing framework is lightweight but powerful. It's built into the Go toolchain, making it easy to get started.

### The `testing` Package

The core of Go testing is the `testing` package. It provides the tools you need to write and run tests.

### The `go test` Command

You run tests using the `go test` command. Here are some common flags:

*   `go test`: Runs all tests in the current directory.
*   `go test ./...`: Runs all tests in the current directory and all subdirectories.
*   `go test -v`: The `-v` flag provides verbose output, showing the name of each test as it runs and whether it passed or failed.
*   `go test -run <regex>`: The `-run` flag allows you to run specific tests that match a regular expression. For example, `go test -run TestLogin` would run only the `TestLogin` test.
*   `go test -cover`: This flag shows test coverage for the packages being tested.

### Naming Conventions

Go's testing tools rely on specific naming conventions:

1.  **File Names:** Test files must end with `_test.go`.
2.  **Test Functions:** Test functions must start with `Test` and take a single argument: `t *testing.T`.

### The `*testing.T` Object

This is your primary tool within a test. It's used to control test flow and report results.

*   `t.Logf(format string, args ...any)`: Prints information to the console if the test fails or if the `-v` flag is used. It does not fail the test.
*   `t.Errorf(format string, args ...any)`: Marks the test as failed but allows it to continue running. This is useful if you want to report multiple errors in a single test.
*   `t.Fatalf(format string, args ...any)`: Marks the test as failed and immediately stops its execution.

## Part 2: A Template for Go Tests

Here is a template that demonstrates a common testing pattern in Go: the table-driven test. This pattern is used extensively in your `gator_cli` project.

```go
package mypackage

import (
	"testing"
)

// MyFunction is the function we want to test.
func MyFunction(input string) string {
	// ... function logic ...
	return "result"
}

// TestMyFunction is the test for MyFunction.
func TestMyFunction(t *testing.T) {
	// Define your test cases as a slice of structs.
	tests := []struct {
		name    string // A descriptive name for the test case
		input   string // The input to pass to the function
		expected string // The expected output
	}{
		// Define individual test cases here.
		{
			name:    "test case 1: a description",
			input:   "some_input",
			expected: "expected_result",
		},
		{
			name:    "test case 2: another description",
			input:   "another_input",
			expected: "another_expected_result",
		},
	}

	// Loop over the test cases.
	for _, tt := range tests {
		// t.Run creates a sub-test, which is great for organization.
		t.Run(tt.name, func(t *testing.T) {
			// Call the function you are testing.
			got := MyFunction(tt.input)

			// Compare the actual result with the expected result.
			if got != tt.expected {
				// If they don't match, report an error.
				t.Errorf("MyFunction() = %v, want %v", got, tt.expected)
			}
		})
	}
}
```

### Explanation of the Template

1.  **`tests := []struct { ... }`**: We define a slice of anonymous structs. Each struct represents a single test case.
2.  **`name`, `input`, `expected`**: These are the fields of our test case struct. `name` is important for identifying which test case failed.
3.  **`for _, tt := range tests`**: We iterate through each test case.
4.  **`t.Run(tt.name, ...)`**: This creates a sub-test. It's good practice because it provides clear output when a specific test case fails.
5.  **`got := MyFunction(tt.input)`**: We execute the function we're testing.
6.  **`if got != tt.expected`**: This is the assertion. We check if the result is what we expected.
7.  **`t.Errorf(...)`**: If the assertion fails, we use `t.Errorf` to report the failure.

## Part 3: Line-by-Line Analysis of `gator_cli` Tests

Here is a detailed breakdown of each test in your project.

### `internal/config/config_test.go`

#### `TestConfig_SetUser`

```go
func TestConfig_SetUser(t *testing.T) {
	// ... (table definition) ...
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			tempFile := filepath.Join(tempDir, configFileName)
			initialData, _ := json.Marshal(tt.config)
			err := os.WriteFile(tempFile, initialData, 0644)
			if err != nil {
				t.Fatalf("failed to create temp config file: %v", err)
			}
			originalGetConfigFilePath := getConfigFilePath
			getConfigFilePath = func() (string, error) {
				return tempFile, nil
			}
			defer func() {
				getConfigFilePath = originalGetConfigFilePath
			}()
			config := tt.config
			err = config.SetUser(tt.username)
			if err != nil {
				t.Errorf("SetUser() error = %v", err)
				return
			}
			if config.CurrentUserName != tt.username {
				t.Errorf("expected CurrentUserName = %q, got %q", tt.username, config.CurrentUserName)
			}
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
*   **Line 5-6**: `tempDir := t.TempDir()` creates a temporary directory. `filepath.Join` constructs a path to a file inside that directory. This ensures the test doesn't affect your actual config file.
*   **Line 7-11**: It sets up the test by writing an initial config file based on the test case data.
*   **Line 12-17**: This is a key part of the test. It replaces the real `getConfigFilePath` function with a fake one that returns the path to our temporary file. `defer` ensures the original function is restored after the test, preventing side effects.
*   **Line 18-23**: It calls `SetUser`, the function being tested, and checks for an error.
*   **Line 24-26**: This is the first assertion. It checks if the `CurrentUserName` in the in-memory `config` object was updated.
*   **Line 27-38**: This is the second assertion. It reads the temporary file from disk, parses it, and checks if the `CurrentUserName` was correctly saved.

#### `TestRead`

This test is very similar to `TestConfig_SetUser` in its setup, using a temporary file and replacing `getConfigFilePath`. The main difference is that it tests the `Read` function.

```go
func TestRead(t *testing.T) {
    // ... (table definition) ...
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            // ... (setup temporary file) ...
			err := tt.setupFile(tempFile)
            // ...
			config, err := Read()

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}
            // ... (assertions) ...
		})
	}
}
```
*   **Line 6**: `tt.setupFile(tempFile)` is a function defined in the test case struct. This allows each test case to set up the temporary file in a different way (e.g., valid JSON, invalid JSON, or no file at all).
*   **Line 8**: `config, err := Read()` calls the function under test.
*   **Line 10-15**: It checks if an error was expected. If so, it verifies that an error was indeed returned. If not, it proceeds to the assertions.

### `cli/commands_test.go`

This file introduces mocks, which are essential for testing code with external dependencies.

#### Mocks (`MockDb`, `MockCfg`)

```go
// MockDb simulates the database
type MockDb struct {
	users       map[string]database.User
	createError error
	resetError  error
}

// MockCfg simulates the config
type MockCfg struct {
	currentUser string
	setUserErr  error
}
```
*   These structs implement the same interfaces as your real database and config handlers. However, their methods are much simpler. They store data in memory (e.g., a `map` for users) and have fields like `createError` that let you tell them to return an error on demand.

#### `TestCommands_Run`

```go
func TestCommands_Run(t *testing.T) {
    // ... (table definition) ...
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commands := Commands{
				CommandMap: make(map[string]func(*State, Command) error),
			}
			tt.setupCmds(&commands)

			mockDb := NewMockDb()
			mockCfg := &MockCfg{}
			state := &State{
				Db:  mockDb,
				Cfg: mockCfg,
			}

			err := commands.Run(state, tt.cmd)

			if tt.expectError {
                // ... (error check) ...
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
			}
		})
	}
}
```
*   **Line 6-8**: A new `Commands` struct is created for each test.
*   **Line 9**: `tt.setupCmds(&commands)` is a function from the test case struct that registers the necessary commands for that specific test.
*   **Line 11-16**: The mocks (`MockDb`, `MockCfg`) are instantiated and placed into a `State` struct. This `state` is what will be passed to the command handlers.
*   **Line 18**: `commands.Run(state, tt.cmd)` executes the command, which will in turn call the handler with the mocked state.
*   **Line 20-27**: The result is checked. The test either asserts that an expected error occurred or that no error occurred.

### `cmd/gator_cli/REPL_test.go`

This file tests the specific command handlers.

#### `TestHandlerLogin`

```go
func TestHandlerLogin(t *testing.T) {
    // ... (table definition) ...
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb := NewMockDb()
			mockCfg := &MockCfg{}
			
			tt.setupDB(mockDb)
			tt.setupCfg(mockCfg)
			
			state := &cli.State{
				Db:  mockDb,
				Cfg: mockCfg,
			}

			err := cli.HandlerLogin(state, tt.cmd)

            // ... (assertions) ...
		})
	}
}
```
*   **Line 6-7**: The mocks are created.
*   **Line 9-10**: The test case's setup functions (`setupDB`, `setupCfg`) are called. This is where you configure the mocks for the specific scenario. For example, for a successful login test, `setupDB` would add a user to the mock database.
*   **Line 12-16**: The `State` is created with the configured mocks.
*   **Line 18**: `cli.HandlerLogin(state, tt.cmd)` calls the handler directly with the mocked state.
*   **The rest**: The assertions check if the handler returned the correct error (or no error) based on the test case.

This same pattern is repeated for `TestHandlerRegister` and `TestHandlerReset`, each time setting up the mocks in a specific way to test the behavior of the handler under different conditions.