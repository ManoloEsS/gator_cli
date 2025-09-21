package main

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ManoloEsS/gator_cli/cli"
	"github.com/ManoloEsS/gator_cli/internal/database"
	"github.com/google/uuid"
)

type MockDb struct {
	users map[string]database.User
}

func NewMockDb() *MockDb {
	return &MockDb{
		users: make(map[string]database.User),
	}
}

func (m *MockDb) GetUser(ctx context.Context, name string) (database.User, error) {
	user, exists := m.users[name]
	if !exists {
		return database.User{}, errors.New("sql: no rows in result set")
	}
	return user, nil
}

func (m *MockDb) CreateUser(ctx context.Context, params database.CreateUserParams) (database.User, error) {
	if _, exists := m.users[params.Name]; exists {
		return database.User{}, errors.New("user already exists")
	}
	user := database.User{
		ID:        params.ID,
		CreatedAt: params.CreatedAt,
		UpdatedAt: params.UpdatedAt,
		Name:      params.Name,
	}
	m.users[params.Name] = user
	return user, nil
}

func (m *MockDb) ResetUsers(ctx context.Context) error {
	m.users = make(map[string]database.User)
	return nil
}

type MockCfg struct {
	currentUser string
	setUserErr  error
}

func (m *MockCfg) SetUser(name string) error {
	if m.setUserErr != nil {
		return m.setUserErr
	}
	m.currentUser = name
	return nil
}

func TestHandlerLogin(t *testing.T) {
	tests := []struct {
		name        string
		cmd         cli.Command
		setupDB     func(*MockDb)
		setupCfg    func(*MockCfg)
		expectError bool
		errorMsg    string
	}{
		{
			name: "login with no arguments",
			cmd: cli.Command{
				Name:      "login",
				Arguments: []string{},
			},
			setupDB:     func(db *MockDb) {},
			setupCfg:    func(cfg *MockCfg) {},
			expectError: true,
			errorMsg:    "usage: login <name>",
		},
		{
			name: "login with non-registered user",
			cmd: cli.Command{
				Name:      "login",
				Arguments: []string{"non-registered-user"},
			},
			setupDB:     func(db *MockDb) {},
			setupCfg:    func(cfg *MockCfg) {},
			expectError: true,
			errorMsg:    "user is not registered in database\n",
		},
		{
			name: "successful login",
			cmd: cli.Command{
				Name:      "login",
				Arguments: []string{"existing-user"},
			},
			setupDB: func(db *MockDb) {
				db.users["existing-user"] = database.User{
					ID:        uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Name:      "existing-user",
				}
			},
			setupCfg:    func(cfg *MockCfg) {},
			expectError: false,
		},
		{
			name: "login with config set user error",
			cmd: cli.Command{
				Name:      "login",
				Arguments: []string{"existing-user"},
			},
			setupDB: func(db *MockDb) {
				db.users["existing-user"] = database.User{
					ID:        uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Name:      "existing-user",
				}
			},
			setupCfg: func(cfg *MockCfg) {
				cfg.setUserErr = errors.New("config error")
			},
			expectError: true,
			errorMsg:    "login handler couldn't switch user: config error\n",
		},
	}

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
		setupDB     func(*MockDb)
		setupCfg    func(*MockCfg)
		expectError bool
		errorMsg    string
	}{
		{
			name: "register with no arguments",
			cmd: cli.Command{
				Name:      "register",
				Arguments: []string{},
			},
			setupDB:     func(db *MockDb) {},
			setupCfg:    func(cfg *MockCfg) {},
			expectError: true,
			errorMsg:    "usage: register <name>\n",
		},
		{
			name: "register existing user",
			cmd: cli.Command{
				Name:      "register",
				Arguments: []string{"existing-user"},
			},
			setupDB: func(db *MockDb) {
				db.users["existing-user"] = database.User{
					ID:        uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Name:      "existing-user",
				}
			},
			setupCfg:    func(cfg *MockCfg) {},
			expectError: true,
			errorMsg:    "user already exists\n",
		},
		{
			name: "successful registration",
			cmd: cli.Command{
				Name:      "register",
				Arguments: []string{"new-user"},
			},
			setupDB:     func(db *MockDb) {},
			setupCfg:    func(cfg *MockCfg) {},
			expectError: false,
		},
		{
			name: "register with config set user error",
			cmd: cli.Command{
				Name:      "register",
				Arguments: []string{"new-user"},
			},
			setupDB: func(db *MockDb) {},
			setupCfg: func(cfg *MockCfg) {
				cfg.setUserErr = errors.New("config error")
			},
			expectError: true,
			errorMsg:    "couldn't set new user in config config error\n",
		},
	}

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
				if _, exists := mockDb.users[tt.cmd.Arguments[0]]; !exists {
					t.Errorf("expected user %s to be created in database", tt.cmd.Arguments[0])
				}
			}
		})
	}
}

func TestHandlerReset(t *testing.T) {
	tests := []struct {
		name    string
		setupDB func(*MockDb)
	}{
		{
			name: "reset empty database",
			setupDB: func(db *MockDb) {},
		},
		{
			name: "reset database with users",
			setupDB: func(db *MockDb) {
				db.users["user1"] = database.User{
					ID:        uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Name:      "user1",
				}
				db.users["user2"] = database.User{
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
			mockDb := NewMockDb()
			mockCfg := &MockCfg{}
			
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
			if len(mockDb.users) != 0 {
				t.Errorf("expected all users to be deleted, but found %d users", len(mockDb.users))
			}
		})
	}
}

