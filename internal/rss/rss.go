package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"
)

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("couldn't create request for RSS feed: %w", err)
	}

	req.Header.Set("User-Agent", "Gator/1.0 (Linux; Custom Client)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("couldn't get a response from the RSS feed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	defer resp.Body.Close()
	rawXML, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("couldn't read response body from RSS feed response: %w", err)
	}

	var feedData RSSFeed
	err = xml.Unmarshal(rawXML, &feedData)
	if err != nil {
		return nil, fmt.Errorf("couldn't unmarshall raw XML data: %w", err)
	}

	feedData.Channel.Title = html.UnescapeString(feedData.Channel.Title)
	feedData.Channel.Description = html.UnescapeString(feedData.Channel.Description)

	for i := range feedData.Channel.Item {
		item := &feedData.Channel.Item[i]
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
	}

	return &feedData, nil
}
