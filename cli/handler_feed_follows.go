package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/ManoloEsS/gator_cli/internal/database"
	"github.com/google/uuid"
)

func HandlerFeedFollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) < 1 {
		return fmt.Errorf("usage: %s <name> <url>\n", cmd.Name)
	}

	url := cmd.Arguments[0]
	feed, err := s.Db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Couldn't retrieve feed data: %w", err)
	}

	feedFollowRow, err := s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Couldn't create feed follow and it to the database: %w", err)
	}

	fmt.Printf("============ Follow added! ============\n")
	fmt.Printf("User: %s\n", feedFollowRow.UserName)
	fmt.Println()
	fmt.Printf("Feed: %s\n", feedFollowRow.FeedName)

	return nil
}

func HandlerFeedFollowsForUser(s *State, cmd Command, user database.User) error {
	feedFollows, err := s.Db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("Couldn't retrieve feed follows for user %s: %w", user.Name, err)
	}
	if len(feedFollows) == 0 {
		fmt.Println("No feed follows found for this user.")
	}

	fmt.Println()
	fmt.Printf("============ User %s's follows ============\n", user.Name)
	fmt.Printf("Feed:               Follow date:\n")
	for _, ff := range feedFollows {
		fmt.Printf("%-20s%02d-%02d-%02d\n", ff.FeedName, ff.CreatedAt.Month(), ff.CreatedAt.Day(), ff.CreatedAt.Year())
	}

	return nil
}

func HandlerUnfollowFeed(s *State, cmd Command, user database.User) error {

	return nil
}
