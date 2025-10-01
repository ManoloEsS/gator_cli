# Testing Documentation

## Overview

This document describes the testing infrastructure and approach for the Gator CLI project.

## Test Structure

### Test Files

- `cli/commands_test.go` - Tests for command registration and execution
- `cli/middleware_test.go` - Tests for authentication middleware
- `cli/handler_user_test.go` - Tests for user management handlers
- `cli/handler_feed_follows_test.go` - Tests for feed follow/unfollow handlers
- `cli/handler_rssfeed_test.go` - Tests for RSS feed handling and HTML stripping
- `cmd/gator_cli/REPL_test.go` - Integration tests for REPL handlers
- `internal/config/config_test.go` - Tests for configuration management
- `internal/rss/rss_test.go` - Tests for RSS feed fetching
- `test/mock_db.go` - Mock database implementation for testing
- `test/mock_config.go` - Mock configuration implementation for testing

## Running Tests

### Run all tests
```bash
go test ./...
```

### Run tests with verbose output
```bash
go test ./... -v
```

### Run tests with coverage
```bash
go test ./... -cover
```

### Run specific test
```bash
go test ./cli -run TestMiddlewareLoggedIn
```

## Test Coverage

Current test coverage:
- `cli`: 45.6%
- `internal/config`: 80.6%
- `internal/rss`: 92.0%

## Mock Infrastructure

### MockDb

The `test.MockDb` struct provides a mock implementation of the `DBInterface`:

- Tracks users, feeds, and feed follows in memory
- Supports error injection for testing error scenarios
- Properly implements all database operations

Example usage:
```go
mockDb := test.NewMockDb()
mockDb.Users["alice"] = database.User{
    ID:   uuid.New(),
    Name: "alice",
}
```

### MockCfg

The `test.MockCfg` struct provides a mock implementation of the `ConfigInterface`:

- Tracks current user
- Supports error injection for SetUser operations

Example usage:
```go
mockCfg := &test.MockCfg{
    CurrentUser: "alice",
}
```

## Testing Best Practices

1. **Use table-driven tests**: Most tests use the table-driven approach for better coverage
2. **Mock external dependencies**: Tests use mock HTTP servers instead of real network calls
3. **Test error cases**: Each handler has tests for both success and error scenarios
4. **Isolate tests**: Each test creates its own mock instances to avoid state pollution
5. **Descriptive test names**: Test names clearly describe what is being tested

## Adding New Tests

When adding new tests:

1. Follow the existing table-driven test pattern
2. Create appropriate test cases for success and error scenarios
3. Use the mock infrastructure provided in `test/` package
4. Ensure tests are isolated and don't depend on each other
5. Run tests with `-v` flag to verify output

## Example Test

```go
func TestHandlerExample(t *testing.T) {
    tests := []struct {
        name        string
        cmd         Command
        setupDB     func(*test.MockDb)
        expectError bool
        errorMsg    string
    }{
        {
            name: "successful operation",
            cmd: Command{
                Name:      "example",
                Arguments: []string{"arg1"},
            },
            setupDB: func(db *test.MockDb) {
                // Setup mock data
            },
            expectError: false,
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockDb := test.NewMockDb()
            mockCfg := &test.MockCfg{}
            tt.setupDB(mockDb)

            state := &State{
                Db:  mockDb,
                Cfg: mockCfg,
            }

            err := HandlerExample(state, tt.cmd)

            if tt.expectError {
                if err == nil {
                    t.Error("expected error but got none")
                }
            } else {
                if err != nil {
                    t.Errorf("expected no error but got: %v", err)
                }
            }
        })
    }
}
```
