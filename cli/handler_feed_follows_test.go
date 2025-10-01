package cli

import (
	"testing"
	"time"

	"github.com/ManoloEsS/gator_cli/internal/database"
	"github.com/ManoloEsS/gator_cli/test"
	"github.com/google/uuid"
)

func TestHandlerFeedFollow(t *testing.T) {
	tests := []struct {
		name        string
		cmd         Command
		user        database.User
		setupDB     func(*test.MockDb)
		expectError bool
		errorMsg    string
	}{
		{
			name: "successful feed follow",
			cmd: Command{
				Name:      "follow",
				Arguments: []string{"https://example.com/feed.xml"},
			},
			user: database.User{
				ID:        uuid.MustParse("12345678-1234-1234-1234-123456789012"),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Name:      "testuser",
			},
			setupDB: func(db *test.MockDb) {
				feedID := uuid.MustParse("22345678-1234-1234-1234-123456789012")
				db.Feeds["https://example.com/feed.xml"] = database.Rssfeed{
					ID:        feedID,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Name:      "Test Feed",
					Url:       "https://example.com/feed.xml",
					UserID:    uuid.New(),
				}
			},
			expectError: false,
		},
		{
			name: "follow with no URL argument",
			cmd: Command{
				Name:      "follow",
				Arguments: []string{},
			},
			user: database.User{
				ID:   uuid.New(),
				Name: "testuser",
			},
			setupDB:     func(db *test.MockDb) {},
			expectError: true,
			errorMsg:    "usage:",
		},
		{
			name: "follow non-existent feed",
			cmd: Command{
				Name:      "follow",
				Arguments: []string{"https://nonexistent.com/feed.xml"},
			},
			user: database.User{
				ID:   uuid.New(),
				Name: "testuser",
			},
			setupDB:     func(db *test.MockDb) {},
			expectError: true,
			errorMsg:    "Couldn't retrieve feed data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb := test.NewMockDb()
			mockCfg := &test.MockCfg{}
			tt.setupDB(mockDb)

			// Add the user to the mock DB
			mockDb.Users[tt.user.Name] = tt.user

			state := &State{
				Db:  mockDb,
				Cfg: mockCfg,
			}

			err := HandlerFeedFollow(state, tt.cmd, tt.user)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				} else if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error to contain '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestHandlerFeedFollowsForUser(t *testing.T) {
	tests := []struct {
		name        string
		cmd         Command
		user        database.User
		setupDB     func(*test.MockDb)
		expectError bool
	}{
		{
			name: "list feed follows for user with follows",
			cmd: Command{
				Name:      "following",
				Arguments: []string{},
			},
			user: database.User{
				ID:        uuid.MustParse("12345678-1234-1234-1234-123456789012"),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Name:      "testuser",
			},
			setupDB: func(db *test.MockDb) {
				userID := uuid.MustParse("12345678-1234-1234-1234-123456789012")
				feedID := uuid.MustParse("22345678-1234-1234-1234-123456789012")
				
				db.FeedFollows[userID] = []database.GetFeedFollowsForUserRow{
					{
						ID:        uuid.New(),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						UserID:    userID,
						FeedID:    feedID,
						FeedName:  "Test Feed",
						Username:  "testuser",
					},
				}
			},
			expectError: false,
		},
		{
			name: "list feed follows for user with no follows",
			cmd: Command{
				Name:      "following",
				Arguments: []string{},
			},
			user: database.User{
				ID:   uuid.New(),
				Name: "testuser",
			},
			setupDB:     func(db *test.MockDb) {},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb := test.NewMockDb()
			mockCfg := &test.MockCfg{}
			tt.setupDB(mockDb)

			state := &State{
				Db:  mockDb,
				Cfg: mockCfg,
			}

			err := HandlerFeedFollowsForUser(state, tt.cmd, tt.user)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestHandlerUnfollowFeed(t *testing.T) {
	tests := []struct {
		name        string
		cmd         Command
		user        database.User
		setupDB     func(*test.MockDb)
		expectError bool
		errorMsg    string
	}{
		{
			name: "successful unfollow",
			cmd: Command{
				Name:      "unfollow",
				Arguments: []string{"https://example.com/feed.xml"},
			},
			user: database.User{
				ID:        uuid.MustParse("12345678-1234-1234-1234-123456789012"),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Name:      "testuser",
			},
			setupDB: func(db *test.MockDb) {
				userID := uuid.MustParse("12345678-1234-1234-1234-123456789012")
				feedID := uuid.MustParse("22345678-1234-1234-1234-123456789012")
				
				db.Feeds["https://example.com/feed.xml"] = database.Rssfeed{
					ID:        feedID,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Name:      "Test Feed",
					Url:       "https://example.com/feed.xml",
					UserID:    uuid.New(),
				}

				db.FeedFollows[userID] = []database.GetFeedFollowsForUserRow{
					{
						ID:        uuid.New(),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						UserID:    userID,
						FeedID:    feedID,
						FeedName:  "Test Feed",
						Username:  "testuser",
					},
				}
			},
			expectError: false,
		},
		{
			name: "unfollow with no URL argument",
			cmd: Command{
				Name:      "unfollow",
				Arguments: []string{},
			},
			user: database.User{
				ID:   uuid.New(),
				Name: "testuser",
			},
			setupDB:     func(db *test.MockDb) {},
			expectError: true,
			errorMsg:    "usage:",
		},
		{
			name: "unfollow non-existent feed",
			cmd: Command{
				Name:      "unfollow",
				Arguments: []string{"https://nonexistent.com/feed.xml"},
			},
			user: database.User{
				ID:   uuid.New(),
				Name: "testuser",
			},
			setupDB:     func(db *test.MockDb) {},
			expectError: true,
			errorMsg:    "Couldn't unfollow feed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb := test.NewMockDb()
			mockCfg := &test.MockCfg{}
			tt.setupDB(mockDb)

			// Add the user to the mock DB
			mockDb.Users[tt.user.Name] = tt.user

			state := &State{
				Db:  mockDb,
				Cfg: mockCfg,
			}

			err := HandlerUnfollowFeed(state, tt.cmd, tt.user)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				} else if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error to contain '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
			}
		})
	}
}
