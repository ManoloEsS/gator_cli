package cli

import (
	"context"
	"fmt"

	"github.com/ManoloEsS/gator_cli/internal/rss"
)

func HandlerAggregate(s *State, cmd Command) error {
	// if len(cmd.Arguments) == 0 {
	// 	return fmt.Errorf("usage: %s <url>\n", cmd.Name)
	// }

	RSSObject, err := rss.FetchFeed(context.Background(), "https://wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	fmt.Printf("%v", RSSObject)

	return nil

}
