package rss

import (
	"context"
	"fmt"
	"testing"
)

// TODO: tests for FetchFeed
func TestFetchFeed(t *testing.T) {
	tests := []struct {
		name         string
		inputContext context.Context
		inputFeed    string
		expected     *RSSFeed
		expectedErr  error
	}{
		{
			name:         "test Fetch Feed output",
			inputContext: context.Background(),
			inputFeed:    "https://wagslane.dev/index.xml",
			expected:     &RSSFeed{},
			expectedErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FetchFeed(tt.inputContext, tt.inputFeed)
			fmt.Printf("%v", got)

			if got != tt.expected {
				t.Errorf("not working %v", err)
			}

			if err != tt.expectedErr {
				t.Errorf("expected nil %v", err)
			}
		})
	}
}
