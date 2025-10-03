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

func TestHandlerLogin(t *testing.T) {
	tests := []struct {
		name        string
		cmd         Command
		setupDB     func(*test.MockDb)
		setupCfg    func(*test.MockCfg)
		expectError bool
	}{
		{
			name: "login existing user",
			cmd: Command{
				Name:      "login",
				Arguments: []string{"alice"},
			},
			setupDB: func(db *test.MockDb) {
				db.Users["alice"] = database.User{
					ID:        uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Name:      "alice",
				}
			},
			setupCfg:    func(cfg *test.MockCfg) {},
			expectError: false,
		},
		{
			name: "login non-existing user",
			cmd: Command{
				Name:      "login",
				Arguments: []string{"bob"},
			},
			setupDB:     func(db *test.MockDb) {},
			setupCfg:    func(cfg *test.MockCfg) {},
			expectError: true,
		},
		{
			name: "login without username",
			cmd: Command{
				Name:      "login",
				Arguments: []string{},
			},
			setupDB:     func(db *test.MockDb) {},
			setupCfg:    func(cfg *test.MockCfg) {},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb := test.NewMockDb()
			mockCfg := &test.MockCfg{}
			tt.setupDB(mockDb)
			tt.setupCfg(mockCfg)

			state := &State{
				Db:  mockDb,
				Cfg: mockCfg,
			}

			err := HandlerLogin(state, tt.cmd)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
				if len(tt.cmd.Arguments) > 0 && mockCfg.CurrentUser != tt.cmd.Arguments[0] {
					t.Errorf("expected current user to be %s but got %s", tt.cmd.Arguments[0], mockCfg.CurrentUser)
				}
			}
		})
	}
}

func TestHandlerRegister(t *testing.T) {
	tests := []struct {
		name        string
		cmd         Command
		setupDB     func(*test.MockDb)
		setupCfg    func(*test.MockCfg)
		expectError bool
	}{
		{
			name: "register new user",
			cmd: Command{
				Name:      "register",
				Arguments: []string{"alice"},
			},
			setupDB:     func(db *test.MockDb) {},
			setupCfg:    func(cfg *test.MockCfg) {},
			expectError: false,
		},
		{
			name: "register existing user",
			cmd: Command{
				Name:      "register",
				Arguments: []string{"alice"},
			},
			setupDB: func(db *test.MockDb) {
				db.Users["alice"] = database.User{
					ID:        uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Name:      "alice",
				}
			},
			setupCfg:    func(cfg *test.MockCfg) {},
			expectError: true,
		},
		{
			name: "register without username",
			cmd: Command{
				Name:      "register",
				Arguments: []string{},
			},
			setupDB:     func(db *test.MockDb) {},
			setupCfg:    func(cfg *test.MockCfg) {},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb := test.NewMockDb()
			mockCfg := &test.MockCfg{}
			tt.setupDB(mockDb)
			tt.setupCfg(mockCfg)

			state := &State{
				Db:  mockDb,
				Cfg: mockCfg,
			}

			err := HandlerRegister(state, tt.cmd)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
				if len(tt.cmd.Arguments) > 0 {
					// Check user was created
					_, exists := mockDb.Users[tt.cmd.Arguments[0]]
					if !exists {
						t.Error("expected user to be created in database")
					}
					// Check user was set as current
					if mockCfg.CurrentUser != tt.cmd.Arguments[0] {
						t.Errorf("expected current user to be %s but got %s", tt.cmd.Arguments[0], mockCfg.CurrentUser)
					}
				}
			}
		})
	}
}
