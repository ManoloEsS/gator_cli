package test

import (
	"errors"
	"testing"
)

func TestMockCfg_SetUser(t *testing.T) {
	tests := []struct {
		name        string
		setupCfg    func(*MockCfg)
		userName    string
		expectError bool
	}{
		{
			name:        "set user successfully",
			setupCfg:    func(cfg *MockCfg) {},
			userName:    "alice",
			expectError: false,
		},
		{
			name: "set user with error",
			setupCfg: func(cfg *MockCfg) {
				cfg.SetUserErr = errors.New("config error")
			},
			userName:    "bob",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCfg := &MockCfg{}
			tt.setupCfg(mockCfg)

			err := mockCfg.SetUser(tt.userName)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
				if mockCfg.CurrentUser != tt.userName {
					t.Errorf("expected CurrentUser to be %s but got %s", tt.userName, mockCfg.CurrentUser)
				}
			}
		})
	}
}

func TestMockCfg_GetCurrentUser(t *testing.T) {
	tests := []struct {
		name         string
		currentUser  string
		expectedUser string
	}{
		{
			name:         "get current user",
			currentUser:  "alice",
			expectedUser: "alice",
		},
		{
			name:         "get empty current user",
			currentUser:  "",
			expectedUser: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCfg := &MockCfg{
				CurrentUser: tt.currentUser,
			}

			user := mockCfg.GetCurrentUser()

			if user != tt.expectedUser {
				t.Errorf("expected user %s but got %s", tt.expectedUser, user)
			}
		})
	}
}
