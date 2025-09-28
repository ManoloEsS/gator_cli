package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/ManoloEsS/gator_cli/internal/database"
	"github.com/google/uuid"
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

	feedFollowRow, err := s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    userData.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Couldn't create feed follow and it to the database: %w", err)
	}

	fmt.Printf("=============Follow added!=============\n")
	fmt.Printf("User: %s\n", feedFollowRow.UserName)
	fmt.Println()
	fmt.Printf("Feed: %s\n", feedFollowRow.FeedName)

	return nil
}

func HandlerFeedFollowsForUser(s *State, cmd Command) error {
	currentUser := s.Cfg.GetCurrentUser()
	currentUserData, err := s.Db.GetUser(context.Background(), currentUser)
	if err != nil {
		return fmt.Errorf("Couldn't retrieve user data from database: %w", err)
	}

	feedFollows, err := s.Db.GetFeedFollowsForUser(context.Background(), currentUserData.ID)
	if err != nil {
		return fmt.Errorf("Couldn't retrieve feed follows for user %s: %w", currentUserData.Name, err)
	}
	if len(feedFollows) == 0 {
		fmt.Println("No feed follows found for this user.")
	}

	fmt.Println()
	fmt.Printf("=============User %s's follows=============\n", currentUserData.Name)
	fmt.Printf("-Feed                -Follow date\n")
	for _, ff := range feedFollows {
		fmt.Printf("> %-20s%02d-%02d-%02d\n", ff.FeedName, ff.CreatedAt.Month(), ff.CreatedAt.Day(), ff.CreatedAt.Year())
		fmt.Println()
	}

	return nil
}
