package cli

import (
	"errors"
	"testing"

	"github.com/ManoloEsS/gator_cli/internal/database"
	"github.com/ManoloEsS/gator_cli/test"
	"github.com/google/uuid"
	"time"
)

func TestMiddlewareLoggedIn(t *testing.T) {
	tests := []struct {
		name            string
		currentUser     string
		setupDB         func(*test.MockDb)
		handlerFunc     func(s *State, cmd Command, user database.User) error
		cmd             Command
		expectError     bool
		expectedErrMsg  string
		validateHandler func(*testing.T, database.User)
	}{
		{
			name:        "successful middleware with valid user",
			currentUser: "testuser",
			setupDB: func(db *test.MockDb) {
				db.Users["testuser"] = database.User{
					ID:        uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Name:      "testuser",
				}
			},
			handlerFunc: func(s *State, cmd Command, user database.User) error {
				// Handler succeeds
				return nil
			},
			cmd:         Command{Name: "test", Arguments: []string{}},
			expectError: false,
			validateHandler: func(t *testing.T, user database.User) {
				if user.Name != "testuser" {
					t.Errorf("expected user.Name = 'testuser', got '%s'", user.Name)
				}
			},
		},
		{
			name:        "user not found in database",
			currentUser: "nonexistent",
			setupDB: func(db *test.MockDb) {
				// Don't add any users
			},
			handlerFunc: func(s *State, cmd Command, user database.User) error {
				t.Error("handler should not be called when user not found")
				return nil
			},
			cmd:            Command{Name: "test", Arguments: []string{}},
			expectError:    true,
			expectedErrMsg: "couldn't retrieve current user's data",
		},
		{
			name:        "handler returns error",
			currentUser: "testuser",
			setupDB: func(db *test.MockDb) {
				db.Users["testuser"] = database.User{
					ID:        uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Name:      "testuser",
				}
			},
			handlerFunc: func(s *State, cmd Command, user database.User) error {
				return errors.New("handler error")
			},
			cmd:            Command{Name: "test", Arguments: []string{}},
			expectError:    true,
			expectedErrMsg: "handler error",
		},
		{
			name:        "middleware passes correct user to handler",
			currentUser: "alice",
			setupDB: func(db *test.MockDb) {
				db.Users["alice"] = database.User{
					ID:        uuid.MustParse("12345678-1234-1234-1234-123456789012"),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Name:      "alice",
				}
			},
			handlerFunc: func(s *State, cmd Command, user database.User) error {
				if user.Name != "alice" {
					return errors.New("wrong user passed to handler")
				}
				if user.ID.String() != "12345678-1234-1234-1234-123456789012" {
					return errors.New("wrong user ID passed to handler")
				}
				return nil
			},
			cmd:         Command{Name: "test", Arguments: []string{"arg1"}},
			expectError: false,
		},
		{
			name:        "empty current user",
			currentUser: "",
			setupDB: func(db *test.MockDb) {
				// No users
			},
			handlerFunc: func(s *State, cmd Command, user database.User) error {
				t.Error("handler should not be called with empty user")
				return nil
			},
			cmd:            Command{Name: "test", Arguments: []string{}},
			expectError:    true,
			expectedErrMsg: "couldn't retrieve current user's data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb := test.NewMockDb()
			mockCfg := &test.MockCfg{
				CurrentUser: tt.currentUser,
			}

			tt.setupDB(mockDb)

			state := &State{
				Db:  mockDb,
				Cfg: mockCfg,
			}

			// Create wrapped handler using middleware
			wrappedHandler := MiddlewareLoggedIn(tt.handlerFunc)

			// Execute the wrapped handler
			err := wrappedHandler(state, tt.cmd)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				} else if tt.expectedErrMsg != "" && !contains(err.Error(), tt.expectedErrMsg) {
					t.Errorf("expected error to contain '%s', got '%s'", tt.expectedErrMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
				if tt.validateHandler != nil {
					// Call validateHandler if there's a user to validate
					if user, exists := mockDb.Users[tt.currentUser]; exists {
						tt.validateHandler(t, user)
					}
				}
			}
		})
	}
}
