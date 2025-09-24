package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/ManoloEsS/gator_cli/internal/database"
	"github.com/google/uuid"
)

func HandlerAddFeed(s *State, cmd Command) error {
	if len(cmd.Arguments) < 2 {
		return fmt.Errorf("usage: %s <name> <url>\n", cmd.Name)
	}

	feedName := cmd.Arguments[0]
	feedUrl := cmd.Arguments[1]
	currentUser := s.Cfg.GetCurrentUser()
	currentUserID, err := s.Db.GetUser(context.Background(), currentUser)
	if err != nil {
		return err
	}

	_, err = s.Db.CreateRSSFeed(context.Background(), database.CreateRSSFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      feedName,
		Url:       feedUrl,
		UserID:    currentUserID.ID,
	})

	fmt.Printf("\"%s succesfully added to %s's feed\n", feedName, currentUser)

	return nil
}
