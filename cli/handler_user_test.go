package cli

import (
	"testing"
	"time"

	"github.com/ManoloEsS/gator_cli/internal/database"
	"github.com/ManoloEsS/gator_cli/test"
	"github.com/google/uuid"
)

func TestHandlerListUsers(t *testing.T) {
	tests := []struct {
		name        string
		cmd         Command
		currentUser string
		setupDB     func(*test.MockDb)
		expectError bool
	}{
		{
			name: "list users with current user marked",
			cmd: Command{
				Name:      "users",
				Arguments: []string{},
			},
			currentUser: "alice",
			setupDB: func(db *test.MockDb) {
				db.Users["alice"] = database.User{
					ID:        uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Name:      "alice",
				}
				db.Users["bob"] = database.User{
					ID:        uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Name:      "bob",
				}
			},
			expectError: false,
		},
		{
			name: "list users with no current user",
			cmd: Command{
				Name:      "users",
				Arguments: []string{},
			},
			currentUser: "",
			setupDB: func(db *test.MockDb) {
				db.Users["alice"] = database.User{
					ID:        uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Name:      "alice",
				}
			},
			expectError: false,
		},
		{
			name: "list users with empty database",
			cmd: Command{
				Name:      "users",
				Arguments: []string{},
			},
			currentUser: "",
			setupDB:     func(db *test.MockDb) {},
			expectError: false,
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

			err := HandlerListUsers(state, tt.cmd)

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

func TestPrintUser(t *testing.T) {
	// Just ensure the function doesn't panic
	user := database.User{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      "testuser",
	}

	// This should not panic
	PrintUser(user)
}
