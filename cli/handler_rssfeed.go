package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ManoloEsS/gator_cli/internal/database"
	"github.com/ManoloEsS/gator_cli/internal/rss"
	"github.com/google/uuid"
)

func HandlerAddFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) < 2 {
		return fmt.Errorf("usage: %s <name> <url>\n", cmd.Name)
	}

	feedName := cmd.Arguments[0]
	feedUrl := cmd.Arguments[1]

	feed, err := s.Db.CreateRSSFeed(context.Background(), database.CreateRSSFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      feedName,
		Url:       feedUrl,
		UserID:    user.ID,
	})

	_, err = s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Couldn't create feed follow from added feed: %w", err)
	}

	fmt.Printf("\"%s\" succesfully added to %s's feed\n", feed.Name, user.Name)

	return nil
}

func HandlerAgg(s *State, cmd Command) error {
	feed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("couldn't fetch feed: %w", err)
	}
	fmt.Printf("Feed: %+v\n", feed)
	return nil
}

func HandlerListFeeds(s *State, cmd Command) error {
	type feeds struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	feedsData, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Couldn't retrieve feeds data from database: %w", err)
	}
	for _, item := range feedsData {
		var decodedData []feeds
		fmt.Println("=============USER=============")
		fmt.Printf(">%s\n", item.UserName)
		fmt.Println("-------------FEEDS-------------")

		err := json.Unmarshal(item.FeedDetails, &decodedData)
		if err != nil {
			fmt.Printf("Failed to decode feed details for user %s: %v\n", item.UserName, err)
			continue
		}

		for _, feed := range decodedData {
			fmt.Printf(">%-20s url:%s\n", feed.Name, feed.Url)
		}
		fmt.Println()
	}
	return nil
}
