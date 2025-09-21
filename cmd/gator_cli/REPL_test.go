package main

import (
	"context"
	"errors"
	"testing"

	"github.com/ManoloEsS/gator_cli/cli"
)

type MockDb struct{}

func (m *MockDb) GetUser(ctx context.Context, name string) (string, error) {
	if name == "non-registered user" {
		return "", errors.New("sql: no rows in result set")
	}
	return "registered-user", nil
}

type MockCfg struct{}

func TestLoginCommand(t *testing.T) {
	cases := []struct {
		name     string
		testArgs []string
		expected string
	}{
		{
			name:     "None-registerd user login",
			testArgs: []string{"login", "non-registered-user"},
			expected: "user is not registered in database\n",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			mockDb := MockDb{}
			mockCfg := MockCfg{}
			state := &cli.State{
				Db:  mockDb,
				Cfg: mockCfg,
			}

			cmd := cli.Command{
				Name:      tt.testArgs[0],
				Arguments: tt.testArgs[1:],
			}

			err := cli.HandlerLogin(state, cmd)

			if err != nil && err.Error() != tt.expected {
				t.Errorf("got error: %q, expected: %q", err.Error(), tt.expected)
			}

			if err == nil && tt.expected != "User login successfull!\n" {
				t.Errorf("expected success message, but got no error\n")
			}
		})
	}
}

