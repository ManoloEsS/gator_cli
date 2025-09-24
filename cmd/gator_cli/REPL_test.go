package main

import (
	"errors"
	"testing"
	"time"

	"github.com/ManoloEsS/gator_cli/cli"
	"github.com/ManoloEsS/gator_cli/internal/database"
	"github.com/ManoloEsS/gator_cli/test"
	"github.com/google/uuid"
)

func TestHandlerLogin(t *testing.T) {
	tests := []struct {
		name        string
		cmd         cli.Command
		setupDB     func(*test.MockDb)
		setupCfg    func(*test.MockCfg)
		expectError bool
		errorMsg    string
	}{
		{
			name: "login with no arguments",
			cmd: cli.Command{
				Name:      "login",
				Arguments: []string{},
			},
			setupDB:     func(db *test.MockDb) {},
			setupCfg:    func(cfg *test.MockCfg) {},
			expectError: true,
			errorMsg:    "usage: login <name>",
		},
		{
			name: "login with non-registered user",
			cmd: cli.Command{
				Name:      "login",
				Arguments: []string{"non-registered-user"},
			},
			setupDB:     func(db *test.MockDb) {},
			setupCfg:    func(cfg *test.MockCfg) {},
			expectError: true,
			errorMsg:    "user is not registered in database\n",
		},
		{
			name: "successful login",
			cmd: cli.Command{
				Name:      "login",
				Arguments: []string{"existing-user"},
			},
			setupDB: func(db *test.MockDb) {
				db.Users["existing-user"] = database.User{
					ID:        uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Name:      "existing-user",
				}
			},
			setupCfg:    func(cfg *test.MockCfg) {},
			expectError: false,
		},
		{
			name: "login with config set user error",
			cmd: cli.Command{
				Name:      "login",
				Arguments: []string{"existing-user"},
			},
			setupDB: func(db *test.MockDb) {
				db.Users["existing-user"] = database.User{
					ID:        uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Name:      "existing-user",
				}
			},
			setupCfg: func(cfg *test.MockCfg) {
				cfg.SetUserErr = errors.New("config error")
			},
			expectError: true,
			errorMsg:    "login handler couldn't switch user: config error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb := test.NewMockDb()
			mockCfg := &test.MockCfg{}

			tt.setupDB(mockDb)
			tt.setupCfg(mockCfg)

			state := &cli.State{
				Db:  mockDb,
				Cfg: mockCfg,
			}

			err := cli.HandlerLogin(state, tt.cmd)

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

func TestHandlerRegister(t *testing.T) {
	tests := []struct {
		name        string
		cmd         cli.Command
		setupDB     func(*test.MockDb)
		setupCfg    func(*test.MockCfg)
		expectError bool
		errorMsg    string
	}{
		{
			name: "register with no arguments",
			cmd: cli.Command{
				Name:      "register",
				Arguments: []string{},
			},
			setupDB:     func(db *test.MockDb) {},
			setupCfg:    func(cfg *test.MockCfg) {},
			expectError: true,
			errorMsg:    "usage: register <name>\n",
		},
		{
			name: "register existing user",
			cmd: cli.Command{
				Name:      "register",
				Arguments: []string{"existing-user"},
			},
			setupDB: func(db *test.MockDb) {
				db.Users["existing-user"] = database.User{
					ID:        uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Name:      "existing-user",
				}
			},
			setupCfg:    func(cfg *test.MockCfg) {},
			expectError: true,
			errorMsg:    "user already exists\n",
		},
		{
			name: "successful registration",
			cmd: cli.Command{
				Name:      "register",
				Arguments: []string{"new-user"},
			},
			setupDB:     func(db *test.MockDb) {},
			setupCfg:    func(cfg *test.MockCfg) {},
			expectError: false,
		},
		{
			name: "register with config set user error",
			cmd: cli.Command{
				Name:      "register",
				Arguments: []string{"new-user"},
			},
			setupDB: func(db *test.MockDb) {},
			setupCfg: func(cfg *test.MockCfg) {
				cfg.SetUserErr = errors.New("config error")
			},
			expectError: true,
			errorMsg:    "couldn't set new user in config config error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb := test.NewMockDb()
			mockCfg := &test.MockCfg{}

			tt.setupDB(mockDb)
			tt.setupCfg(mockCfg)

			state := &cli.State{
				Db:  mockDb,
				Cfg: mockCfg,
			}

			err := cli.HandlerRegister(state, tt.cmd)

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
				// Verify user was created in database
				if _, exists := mockDb.Users[tt.cmd.Arguments[0]]; !exists {
					t.Errorf("expected user %s to be created in database", tt.cmd.Arguments[0])
				}
			}
		})
	}
}

func TestHandlerReset(t *testing.T) {
	tests := []struct {
		name    string
		setupDB func(*test.MockDb)
	}{
		{
			name:    "reset empty database",
			setupDB: func(db *test.MockDb) {},
		},
		{
			name: "reset database with users",
			setupDB: func(db *test.MockDb) {
				db.Users["user1"] = database.User{
					ID:        uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Name:      "user1",
				}
				db.Users["user2"] = database.User{
					ID:        uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Name:      "user2",
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb := test.NewMockDb()
			mockCfg := &test.MockCfg{}

			tt.setupDB(mockDb)

			state := &cli.State{
				Db:  mockDb,
				Cfg: mockCfg,
			}

			cmd := cli.Command{
				Name:      "reset",
				Arguments: []string{},
			}

			err := cli.HandlerReset(state, cmd)

			if err != nil {
				t.Errorf("expected no error but got: %v", err)
			}

			// Verify all users were deleted
			if len(mockDb.Users) != 0 {
				t.Errorf("expected all users to be deleted, but found %d users", len(mockDb.Users))
			}
		})
	}
}
