package rss

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFetchFeed(t *testing.T) {
	tests := []struct {
		name           string
		inputContext   context.Context
		mockResponse   string
		mockStatusCode int
		expectError    bool
		validateResult func(*testing.T, *RSSFeed)
	}{
		{
			name:           "valid RSS feed",
			inputContext:   context.Background(),
			mockStatusCode: http.StatusOK,
			mockResponse: `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
	<channel>
		<title>Test Feed</title>
		<link>https://example.com</link>
		<description>Test Description</description>
		<item>
			<title>Test Item</title>
			<link>https://example.com/item1</link>
			<description>Test item description</description>
			<pubDate>Mon, 02 Jan 2006 15:04:05 MST</pubDate>
		</item>
	</channel>
</rss>`,
			expectError: false,
			validateResult: func(t *testing.T, feed *RSSFeed) {
				if feed == nil {
					t.Error("expected feed to not be nil")
					return
				}
				if feed.Channel.Title != "Test Feed" {
					t.Errorf("expected title 'Test Feed', got '%s'", feed.Channel.Title)
				}
				if feed.Channel.Link != "https://example.com" {
					t.Errorf("expected link 'https://example.com', got '%s'", feed.Channel.Link)
				}
				if feed.Channel.Description != "Test Description" {
					t.Errorf("expected description 'Test Description', got '%s'", feed.Channel.Description)
				}
				if len(feed.Channel.Item) != 1 {
					t.Errorf("expected 1 item, got %d", len(feed.Channel.Item))
					return
				}
				if feed.Channel.Item[0].Title != "Test Item" {
					t.Errorf("expected item title 'Test Item', got '%s'", feed.Channel.Item[0].Title)
				}
			},
		},
		{
			name:           "HTML entities are unescaped",
			inputContext:   context.Background(),
			mockStatusCode: http.StatusOK,
			mockResponse: `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
	<channel>
		<title>Test &amp; Feed</title>
		<link>https://example.com</link>
		<description>Description with &lt;tags&gt;</description>
		<item>
			<title>Item &amp; Title</title>
			<link>https://example.com/item1</link>
			<description>Description with &quot;quotes&quot;</description>
			<pubDate>Mon, 02 Jan 2006 15:04:05 MST</pubDate>
		</item>
	</channel>
</rss>`,
			expectError: false,
			validateResult: func(t *testing.T, feed *RSSFeed) {
				if feed == nil {
					t.Error("expected feed to not be nil")
					return
				}
				if feed.Channel.Title != "Test & Feed" {
					t.Errorf("expected title 'Test & Feed', got '%s'", feed.Channel.Title)
				}
				if feed.Channel.Description != "Description with <tags>" {
					t.Errorf("expected description 'Description with <tags>', got '%s'", feed.Channel.Description)
				}
				if len(feed.Channel.Item) > 0 {
					if feed.Channel.Item[0].Title != "Item & Title" {
						t.Errorf("expected item title 'Item & Title', got '%s'", feed.Channel.Item[0].Title)
					}
					if feed.Channel.Item[0].Description != `Description with "quotes"` {
						t.Errorf("expected item description with unescaped quotes, got '%s'", feed.Channel.Item[0].Description)
					}
				}
			},
		},
		{
			name:           "empty feed",
			inputContext:   context.Background(),
			mockStatusCode: http.StatusOK,
			mockResponse: `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
	<channel>
		<title>Empty Feed</title>
		<link>https://example.com</link>
		<description>No items</description>
	</channel>
</rss>`,
			expectError: false,
			validateResult: func(t *testing.T, feed *RSSFeed) {
				if feed == nil {
					t.Error("expected feed to not be nil")
					return
				}
				if len(feed.Channel.Item) != 0 {
					t.Errorf("expected 0 items, got %d", len(feed.Channel.Item))
				}
			},
		},
		{
			name:           "HTTP error",
			inputContext:   context.Background(),
			mockStatusCode: http.StatusNotFound,
			mockResponse:   "Not Found",
			expectError:    true,
		},
		{
			name:           "invalid XML",
			inputContext:   context.Background(),
			mockStatusCode: http.StatusOK,
			mockResponse:   "this is not valid XML",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Check User-Agent header
				if userAgent := r.Header.Get("User-Agent"); !strings.HasPrefix(userAgent, "Gator/") {
					t.Errorf("expected User-Agent to start with 'Gator/', got '%s'", userAgent)
				}

				w.WriteHeader(tt.mockStatusCode)
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			// Call FetchFeed with the test server URL
			got, err := FetchFeed(tt.inputContext, server.URL)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
				if tt.validateResult != nil {
					tt.validateResult(t, got)
				}
			}
		})
	}
}

func TestFetchFeedContextCancellation(t *testing.T) {
	// Create a server that never responds
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Never respond
		select {}
	}))
	defer server.Close()

	// Create a context that is already canceled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := FetchFeed(ctx, server.URL)
	if err == nil {
		t.Error("expected error from canceled context but got none")
	}
}

