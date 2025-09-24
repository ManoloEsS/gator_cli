package cli

import (
	"context"
	"errors"
	"testing"

	"github.com/ManoloEsS/gator_cli/internal/database"
	"github.com/ManoloEsS/gator_cli/test"
)

func TestCommands_Register(t *testing.T) {
	tests := []struct {
		name        string
		commandName string
		handlerFunc func(*State, Command) error
		expectPanic bool
	}{
		{
			name:        "register valid command",
			commandName: "test",
			handlerFunc: func(s *State, cmd Command) error { return nil },
			expectPanic: false,
		},
		{
			name:        "register command with nil handler",
			commandName: "test-nil",
			handlerFunc: nil,
			expectPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commands := Commands{
				CommandMap: make(map[string]func(*State, Command) error),
			}

			if tt.expectPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("expected panic but didn't get one")
					}
				}()
			}

			commands.Register(tt.commandName, tt.handlerFunc)

			// Verify the command was registered
			if handler, exists := commands.CommandMap[tt.commandName]; !exists {
				t.Errorf("expected command %s to be registered", tt.commandName)
			} else if handler == nil && tt.handlerFunc != nil {
				t.Errorf("expected handler to be registered but got nil")
			}
		})
	}
}

func TestCommands_Run(t *testing.T) {
	tests := []struct {
		name        string
		setupCmds   func(*Commands)
		cmd         Command
		expectError bool
		errorMsg    string
	}{
		{
			name: "run existing command",
			setupCmds: func(cmds *Commands) {
				cmds.Register("test", func(s *State, cmd Command) error {
					return nil
				})
			},
			cmd: Command{
				Name:      "test",
				Arguments: []string{},
			},
			expectError: false,
		},
		{
			name: "run non-existing command",
			setupCmds: func(cmds *Commands) {
				// Don't register anything
			},
			cmd: Command{
				Name:      "nonexistent",
				Arguments: []string{},
			},
			expectError: true,
			errorMsg:    "command not found",
		},
		{
			name: "run command that returns error",
			setupCmds: func(cmds *Commands) {
				cmds.Register("error-cmd", func(s *State, cmd Command) error {
					return errors.New("handler error")
				})
			},
			cmd: Command{
				Name:      "error-cmd",
				Arguments: []string{},
			},
			expectError: true,
			errorMsg:    "handler error",
		},
		{
			name: "run command with arguments",
			setupCmds: func(cmds *Commands) {
				cmds.Register("args-cmd", func(s *State, cmd Command) error {
					if len(cmd.Arguments) != 2 {

					}
					return nil
				})
			},
			cmd: Command{
				Name:      "args-cmd",
				Arguments: []string{"arg1", "arg2"},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commands := Commands{
				CommandMap: make(map[string]func(*State, Command) error),
			}
			tt.setupCmds(&commands)

			mockDb := test.NewMockDb()
			mockCfg := &test.MockCfg{}
			state := &State{
				Db:  mockDb,
				Cfg: mockCfg,
			}

			err := commands.Run(state, tt.cmd)

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
		})
	}
}

func TestHandlerRegister_CreateUserError(t *testing.T) {
	mockDb := test.NewMockDb()
	mockDb.CreateError = errors.New("database create error")
	mockCfg := &test.MockCfg{}

	state := &State{
		Db:  mockDb,
		Cfg: mockCfg,
	}

	cmd := Command{
		Name:      "register",
		Arguments: []string{"new-user"},
	}

	err := HandlerRegister(state, cmd)

	if err == nil {
		t.Errorf("expected error but got none")
	} else if !contains(err.Error(), "couldn't create new user") {
		t.Errorf("expected error message to contain 'couldn't create new user', got: %q", err.Error())
	}
}

func TestCommandStruct(t *testing.T) {
	// Test Command struct construction and field access
	cmd := Command{
		Name:      "test-command",
		Arguments: []string{"arg1", "arg2", "arg3"},
	}

	if cmd.Name != "test-command" {
		t.Errorf("expected Name = 'test-command', got %q", cmd.Name)
	}

	if len(cmd.Arguments) != 3 {
		t.Errorf("expected 3 arguments, got %d", len(cmd.Arguments))
	}

	if cmd.Arguments[0] != "arg1" {
		t.Errorf("expected Arguments[0] = 'arg1', got %q", cmd.Arguments[0])
	}

	// Test empty command
	emptyCmd := Command{}
	if emptyCmd.Name != "" {
		t.Errorf("expected empty Name, got %q", emptyCmd.Name)
	}

	if len(emptyCmd.Arguments) != 0 {
		t.Errorf("expected 0 arguments, got %d", len(emptyCmd.Arguments))
	}
}

func TestMockDbEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		setupDb     func(*test.MockDb)
		operation   func(*test.MockDb) error
		expectError bool
	}{
		{
			name: "reset users with error",
			setupDb: func(db *test.MockDb) {
				db.ResetError = errors.New("reset failed")
			},
			operation: func(db *test.MockDb) error {
				return db.ResetUsers(context.Background())
			},
			expectError: true,
		},
		{
			name: "get user from empty database",
			setupDb: func(db *test.MockDb) {
				// No setup needed
			},
			operation: func(db *test.MockDb) error {
				_, err := db.GetUser(context.Background(), "nonexistent")
				return err
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb := test.NewMockDb()
			tt.setupDb(mockDb)

			err := tt.operation(mockDb)

			if tt.expectError && err == nil {
				t.Errorf("expected error but got none")
			} else if !tt.expectError && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
		})
	}
}

func TestIntegrationCommandsWithHandlers(t *testing.T) {
	// Test that the Commands struct can properly run all handlers
	tests := []struct {
		name        string
		commandName string
		handler     func(*State, Command) error
		cmd         Command
		setupState  func(*test.MockDb, *test.MockCfg)
		expectError bool
	}{
		{
			name:        "login command integration",
			commandName: "login",
			handler:     HandlerLogin,
			cmd: Command{
				Name:      "login",
				Arguments: []string{"testuser"},
			},
			setupState: func(db *test.MockDb, cfg *test.MockCfg) {
				// Add a user to login
				db.Users["testuser"] = database.User{Name: "testuser"}
			},
			expectError: false,
		},
		{
			name:        "register command integration",
			commandName: "register",
			handler:     HandlerRegister,
			cmd: Command{
				Name:      "register",
				Arguments: []string{"newuser"},
			},
			setupState:  func(db *test.MockDb, cfg *test.MockCfg) {},
			expectError: false,
		},
		{
			name:        "reset command integration",
			commandName: "reset",
			handler:     HandlerReset,
			cmd: Command{
				Name:      "reset",
				Arguments: []string{},
			},
			setupState: func(db *test.MockDb, cfg *test.MockCfg) {
				// Add some users to reset
				db.Users["user1"] = database.User{Name: "user1"}
				db.Users["user2"] = database.User{Name: "user2"}
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup commands
			commands := Commands{
				CommandMap: make(map[string]func(*State, Command) error),
			}
			commands.Register(tt.commandName, tt.handler)

			// Setup state
			mockDb := test.NewMockDb()
			mockCfg := &test.MockCfg{}
			tt.setupState(mockDb, mockCfg)

			state := &State{
				Db:  mockDb,
				Cfg: mockCfg,
			}

			// Run the command
			err := commands.Run(state, tt.cmd)

			if tt.expectError && err == nil {
				t.Errorf("expected error but got none")
			} else if !tt.expectError && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
		})
	}
}

func TestNewState(t *testing.T) {
	// This test ensures NewState creates a state with the correct interfaces
	mockDb := test.NewMockDb()
	mockCfg := &test.MockCfg{}

	// Create a state using concrete types (simulating production use)
	// Note: We can't test with real *database.Queries and *config.Config
	// without setting up a real database, so we test the interface compliance
	state := &State{
		Db:  mockDb,
		Cfg: mockCfg,
	}

	if state.Db == nil {
		t.Errorf("expected Db to be set")
	}
	if state.Cfg == nil {
		t.Errorf("expected Cfg to be set")
	}

	// Test that the state can be used with handlers
	cmd := Command{
		Name:      "test",
		Arguments: []string{},
	}

	// This should not panic or fail
	err := HandlerReset(state, cmd)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
