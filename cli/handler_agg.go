package cli

import (
	"context"
	"fmt"

	"github.com/ManoloEsS/gator_cli/internal/rss"
)

func HandlerAgg(s *State, cmd Command) error {
	feed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("couldn't fetch feed: %w", err)
	}
	fmt.Printf("Feed: %+v\n", feed)
	return nil
}
