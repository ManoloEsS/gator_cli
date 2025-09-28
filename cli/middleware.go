package cli

import (
	"context"
	"fmt"

	"github.com/ManoloEsS/gator_cli/internal/database"
)

func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		user := s.Cfg.GetCurrentUser()
		userData, err := s.Db.GetUser(context.Background(), user)
		if err != nil {
			return fmt.Errorf("couldn't retrieve current user's data: %w", err)
		}
		return handler(s, cmd, userData)
	}
}
