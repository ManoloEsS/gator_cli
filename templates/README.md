# Test Templates for gator_cli

This directory contains test file templates designed for educational purposes. These templates demonstrate Go testing best practices and provide a structured approach to learning how to write comprehensive tests.

## Overview

The templates are based on the existing test files in the `gator_cli` project and showcase various testing patterns and techniques used in Go development.

## Templates Included

### 1. `config_test_template.go`
**Purpose**: Template for testing configuration file operations (reading, writing, user management)

**Key Learning Points**:
- Table-driven test pattern
- Temporary file handling for filesystem operations  
- Function override patterns for testability
- JSON marshalling/unmarshalling testing
- Setup and cleanup using `t.TempDir()` and `defer`

**Functions Covered**:
- `TestConfig_SetUser` - Testing user configuration updates
- `TestRead` - Testing configuration file reading with various scenarios
- `TestWrite` - Testing configuration file writing
- `TestGetConfigFilePath` - Testing path construction logic

### 2. `commands_test_template.go`
**Purpose**: Template for testing CLI command registration, execution, and handler functions

**Key Learning Points**:
- Interface-based testing with mocks
- Command argument validation
- Error handling and message verification
- Integration testing between components
- Mock setup and state verification

**Functions Covered**:
- `TestCommands_Register` - Testing command registration
- `TestCommands_Run` - Testing command execution
- `TestHandlerRegister` - Testing user registration handler
- `TestHandlerLogin` - Testing user login handler
- `TestHandlerReset` - Testing database reset handler
- `TestIntegrationCommandsWithHandlers` - Integration testing

### 3. `REPL_test_template.go`
**Purpose**: Template for testing REPL (Read-Eval-Print Loop) command handlers

**Key Learning Points**:
- Handler function testing with mock state
- Database interaction testing
- Configuration state management
- Success and error path verification
- Command argument processing

**Functions Covered**:
- `TestHandlerLogin` - Testing login command handler
- `TestHandlerRegister` - Testing register command handler  
- `TestHandlerReset` - Testing reset command handler

### 4. `rss_test_template.go`
**Purpose**: Template for testing RSS feed fetching, parsing, and processing

**Key Learning Points**:
- HTTP client mocking for external API testing
- Context handling in network operations
- XML/RSS parsing validation
- Timeout and cancellation testing
- Test server setup for HTTP testing

**Functions Covered**:
- `TestFetchFeed` - Testing RSS feed fetching from URLs
- `TestRSSFeedParsing` - Testing RSS XML parsing logic
- `TestRSSItemValidation` - Testing RSS item validation
- `TestRSSDateParsing` - Testing RSS date format parsing

### 5. `mock_test_utilities_template.go`
**Purpose**: Template for creating comprehensive mock implementations for testing

**Key Learning Points**:
- Interface-based mocking for clean testing
- Configurable mock behavior for different scenarios
- State tracking and verification in tests
- Error injection for testing error paths
- Builder pattern for complex mock setup

**Components Provided**:
- `MockDb` - Database interface mock
- `MockCfg` - Configuration interface mock
- Helper functions for common test scenarios
- Builder pattern for fluent mock setup
- Assertion helpers for common verifications

## How to Use These Templates

### Step 1: Choose the Appropriate Template
Select the template that matches the functionality you want to test:
- Configuration operations → `config_test_template.go`
- CLI commands and handlers → `commands_test_template.go` or `REPL_test_template.go`
- RSS/HTTP operations → `rss_test_template.go`
- Need mocks → `mock_test_utilities_template.go`

### Step 2: Copy and Customize
1. Copy the template to your test directory
2. Rename the file to match your package (e.g., `your_package_test.go`)
3. Update the package declaration at the top
4. Update import paths to match your project structure

### Step 3: Replace TODO Comments
Each template contains `TODO` comments indicating where you need to add your specific logic:

```go
// TODO: Add test case for setting user on empty config
name:     "TODO: add descriptive name",
username: "TODO: add test username", 
config:   Config{}, // TODO: set initial config state
```

### Step 4: Implement Test Logic
Follow the patterns shown in the templates:

```go
// Example of implementing a test case
{
    name: "set user on empty config",
    username: "testuser",
    config: Config{},
},
```

### Step 5: Add Your Test Cases
Extend the test tables with additional scenarios relevant to your code:

```go
tests := []struct {
    name     string
    username string 
    config   Config
    expectError bool  // Add fields as needed
}{
    // TODO sections show you where to add cases
    {
        name: "your specific test case",
        username: "your test data",
        config: Config{/* your config */},
    },
}
```

### Step 6: Run and Iterate
Run your tests and refine them:

```bash
go test -v ./path/to/your/package
```

## Testing Best Practices Demonstrated

### Table-Driven Tests
All templates use the table-driven test pattern:

```go
tests := []struct {
    name        string
    input       InputType
    expected    ExpectedType
    expectError bool
}{ /* test cases */ }

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // Test implementation
    })
}
```

### Mock Usage
Templates show how to create and use mocks effectively:

```go
mockDb := NewMockDb()
mockDb.Users["testuser"] = database.User{Name: "testuser"}
state := &State{Db: mockDb, Cfg: mockCfg}
```

### Error Testing
Templates demonstrate proper error testing:

```go
if tt.expectError {
    if err == nil {
        t.Errorf("expected error but got none")
    } else if err.Error() != tt.errorMsg {
        t.Errorf("expected error: %q, got: %q", tt.errorMsg, err.Error())
    }
} else {
    if err != nil {
        t.Errorf("expected no error but got: %v", err)
    }
}
```

### Resource Cleanup
Templates show proper cleanup patterns:

```go
tempDir := t.TempDir()  // Automatic cleanup
defer func() {
    // Cleanup code
}()
```

## Common Testing Patterns

### 1. Testing File Operations
Use `t.TempDir()` for temporary files and override global functions for testability.

### 2. Testing HTTP Operations  
Use `httptest.Server` to create test HTTP servers.

### 3. Testing Database Operations
Use mocks that implement the same interfaces as your real database layer.

### 4. Testing Configuration
Override file paths and use temporary directories.

### 5. Testing Command Handlers
Use mocks for dependencies and verify both success and error cases.

## Running Tests

To run tests created from these templates:

```bash
# Run all tests
go test -v ./...

# Run specific package tests
go test -v ./internal/config

# Run with coverage
go test -cover ./...

# Run specific test function
go test -v -run TestFunctionName ./package
```

## Learning Progression

1. **Start with simple tests** - Begin with basic functionality tests
2. **Add error cases** - Test error handling and edge cases
3. **Use mocks** - Learn interface-based testing with mocks
4. **Add integration tests** - Test component interactions
5. **Optimize and refactor** - Improve test maintainability

## Additional Resources

- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Table-Driven Tests in Go](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)
- [Go Testing Best Practices](https://github.com/golang/go/wiki/TestingPolicies)

## Questions and Learning

As you work through these templates, consider:

1. What edge cases might break your function?
2. How can you test error conditions?
3. What dependencies does your code have?
4. How can you make your code more testable?
5. What would happen if external services are unavailable?

Remember: The goal is to learn by doing. Start simple and gradually add complexity as you become more comfortable with Go testing patterns.