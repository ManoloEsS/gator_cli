package cli

import (
	"context"
	"fmt"
)

func HandlerFeedFollow(s *State, cmd Command) error {
	if len(cmd.Arguments) < 1 {
		return fmt.Errorf("usage: %s <name> <url>\n", cmd.Name)
	}

	url := cmd.Arguments[0]
	feed, err := s.Db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Couldn't retrieve feed data: %w", err)
	}

	currentUser := s.Cfg.GetCurrentUser()
	userData, err := s.Db.GetUser(context.Background(), currentUser)
	if err != nil {
		return fmt.Errorf("Couldn't retrieve user data: %w", err)
	}

	return nil
}
