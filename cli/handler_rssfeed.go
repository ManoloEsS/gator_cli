package cli

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/ManoloEsS/gator_cli/internal/database"
	"github.com/ManoloEsS/gator_cli/internal/rss"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/net/html"
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
	if len(cmd.Arguments) < 1 {
		return fmt.Errorf("usage %s <duration>", cmd.Name)
	}
	time_between_reqs := cmd.Arguments[0]
	duration, err := time.ParseDuration(time_between_reqs)
	if err != nil {
		return fmt.Errorf("usage eg: 1s (s: second, m: minute, h: hour): %w", err)
	}
	ticker := time.NewTicker(duration)
	fmt.Printf("Collecting feeds every %s\n", time_between_reqs)

	for ; ; <-ticker.C {
		ScrapeFeeds(s)
	}
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

func HandlerBrowse(s *State, cmd Command, user database.User) error {
	var postsNum int32 = 2
	if len(cmd.Arguments) == 1 {
		if n, err := strconv.Atoi(cmd.Arguments[0]); err == nil {
			postsNum = int32(n)
		}

	}
	posts, err := s.Db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  postsNum,
	})
	if err != nil {
		log.Printf("Couldn't get post for user: %v", err)
	}

	for _, item := range posts {
		fmt.Printf("%s from %s\n", item.PublishedAt.Time.Format("Mon Jan 2"), item.FeedName)
		fmt.Printf("---%s---\n", item.Title)
		fmt.Printf("   %v\n", StripHTML(item.Description.String))
		fmt.Printf("Link: %s\n", item.Url)
		fmt.Println("====================================")
	}
	return nil
}

func scrapeFeed(db DBInterface, feed database.Rssfeed) {
	err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("couldn't mark feed %s fetched: %v", feed.Name, err)
		return
	}

	rssResponseData, err := rss.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("couldn't fetch from feed %s: %w", feed.Url, err)
		return
	}

	for _, item := range rssResponseData.Channel.Item {
		publishedDateParsed := sql.NullTime{}
		if t, ok := parsePubDate(item.PubDate); ok {
			publishedDateParsed = newNullTime(t)
		}

		_, err := db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title:     item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			PublishedAt: publishedDateParsed,
			FeedID:      feed.ID,
		})
		if err != nil {
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
				continue
			}
			log.Printf("couldn't add post to database: %v", err)
			continue
		}
	}

	fmt.Println("===============================================")
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssResponseData.Channel.Item))
}
func ScrapeFeeds(s *State) {
	nextFeed, err := s.Db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Println("couldn't get next feed to fetch", err)
		return
	}
	log.Println("Found a feed to fetch!")
	scrapeFeed(s.Db, nextFeed)
}

func newNullTime(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{}
	}
	return sql.NullTime{
		Time:  t,
		Valid: true,
	}
}

var rssLayouts = []string{
	time.RFC1123,                      // "Mon, 02 Jan 2006 15:04:05 MST"
	time.RFC1123Z,                     // "Mon, 02 Jan 2006 15:04:05 -0700"
	time.RFC822,                       // "02 Jan 06 15:04 MST"
	time.RFC822Z,                      // "02 Jan 06 15:04 -0700"
	time.RFC3339,                      // "2006-01-02T15:04:05Z07:00"
	"Mon, 02 Jan 2006 15:04 MST",      // no seconds
	"Mon, 02 Jan 2006 15:04:05 -0700", // explicit offset
}

func parsePubDate(s string) (time.Time, bool) {
	s = strings.TrimSpace(s)
	for _, layout := range rssLayouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

func StripHTML(s string) string {
	z := html.NewTokenizer(strings.NewReader(s))
	var out strings.Builder
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return strings.TrimSpace(out.String())
		case html.TextToken:
			out.WriteString(string(z.Text()))
		}
	}
}
