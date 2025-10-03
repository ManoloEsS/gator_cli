package test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ManoloEsS/gator_cli/internal/database"
	"github.com/google/uuid"
)

func TestNewMockDb(t *testing.T) {
	mockDb := NewMockDb()
	if mockDb == nil {
		t.Fatal("expected NewMockDb to return non-nil MockDb")
	}
	if mockDb.Users == nil {
		t.Error("expected Users map to be initialized")
	}
	if mockDb.Feeds == nil {
		t.Error("expected Feeds map to be initialized")
	}
	if mockDb.FeedFollows == nil {
		t.Error("expected FeedFollows map to be initialized")
	}
}

func TestMockDb_GetUser(t *testing.T) {
	tests := []struct {
		name        string
		setupDB     func(*MockDb)
		userName    string
		expectError bool
	}{
		{
			name: "get existing user",
			setupDB: func(db *MockDb) {
				db.Users["alice"] = database.User{
					ID:        uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Name:      "alice",
				}
			},
			userName:    "alice",
			expectError: false,
		},
		{
			name:        "get non-existing user",
			setupDB:     func(db *MockDb) {},
			userName:    "bob",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb := NewMockDb()
			tt.setupDB(mockDb)

			user, err := mockDb.GetUser(context.Background(), tt.userName)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
				if user.Name != tt.userName {
					t.Errorf("expected user name %s but got %s", tt.userName, user.Name)
				}
			}
		})
	}
}

func TestMockDb_CreateUser(t *testing.T) {
	tests := []struct {
		name        string
		setupDB     func(*MockDb)
		userName    string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "create new user",
			setupDB:     func(db *MockDb) {},
			userName:    "alice",
			expectError: false,
		},
		{
			name: "create duplicate user",
			setupDB: func(db *MockDb) {
				db.Users["alice"] = database.User{
					ID:        uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Name:      "alice",
				}
			},
			userName:    "alice",
			expectError: true,
			errorMsg:    "user already exists",
		},
		{
			name: "create user with CreateError set",
			setupDB: func(db *MockDb) {
				db.CreateError = errors.New("database error")
			},
			userName:    "bob",
			expectError: true,
			errorMsg:    "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb := NewMockDb()
			tt.setupDB(mockDb)

			params := database.CreateUserParams{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Name:      tt.userName,
			}

			user, err := mockDb.CreateUser(context.Background(), params)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
				if user.Name != tt.userName {
					t.Errorf("expected user name %s but got %s", tt.userName, user.Name)
				}
				// Check if user was added to the map
				if _, exists := mockDb.Users[tt.userName]; !exists {
					t.Error("expected user to be added to Users map")
				}
			}
		})
	}
}

func TestMockDb_ResetUsers(t *testing.T) {
	tests := []struct {
		name        string
		setupDB     func(*MockDb)
		expectError bool
	}{
		{
			name: "reset users with existing users",
			setupDB: func(db *MockDb) {
				db.Users["alice"] = database.User{
					ID:        uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Name:      "alice",
				}
				db.Users["bob"] = database.User{
					ID:        uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Name:      "bob",
				}
			},
			expectError: false,
		},
		{
			name: "reset users with ResetError set",
			setupDB: func(db *MockDb) {
				db.ResetError = errors.New("reset error")
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb := NewMockDb()
			tt.setupDB(mockDb)

			err := mockDb.ResetUsers(context.Background())

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
				if len(mockDb.Users) != 0 {
					t.Errorf("expected Users map to be empty, but has %d users", len(mockDb.Users))
				}
			}
		})
	}
}

func TestMockDb_GetUsers(t *testing.T) {
	tests := []struct {
		name          string
		setupDB       func(*MockDb)
		expectedCount int
	}{
		{
			name:          "get users with empty database",
			setupDB:       func(db *MockDb) {},
			expectedCount: 0,
		},
		{
			name: "get users with multiple users",
			setupDB: func(db *MockDb) {
				db.Users["alice"] = database.User{
					ID:        uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Name:      "alice",
				}
				db.Users["bob"] = database.User{
					ID:        uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Name:      "bob",
				}
			},
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb := NewMockDb()
			tt.setupDB(mockDb)

			users, err := mockDb.GetUsers(context.Background())

			if err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
			if len(users) != tt.expectedCount {
				t.Errorf("expected %d users but got %d", tt.expectedCount, len(users))
			}
		})
	}
}

func TestMockDb_GetFeedByUrl(t *testing.T) {
	tests := []struct {
		name        string
		setupDB     func(*MockDb)
		url         string
		expectError bool
	}{
		{
			name: "get existing feed",
			setupDB: func(db *MockDb) {
				db.Feeds["https://example.com/feed"] = database.Rssfeed{
					ID:        uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Name:      "Example Feed",
					Url:       "https://example.com/feed",
					UserID:    uuid.New(),
				}
			},
			url:         "https://example.com/feed",
			expectError: false,
		},
		{
			name:        "get non-existing feed",
			setupDB:     func(db *MockDb) {},
			url:         "https://nonexistent.com/feed",
			expectError: true,
		},
		{
			name: "get feed with GetFeedError set",
			setupDB: func(db *MockDb) {
				db.GetFeedError = errors.New("database error")
			},
			url:         "https://example.com/feed",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb := NewMockDb()
			tt.setupDB(mockDb)

			feed, err := mockDb.GetFeedByUrl(context.Background(), tt.url)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
				if feed.Url != tt.url {
					t.Errorf("expected feed url %s but got %s", tt.url, feed.Url)
				}
			}
		})
	}
}

func TestMockDb_CreateFeedFollow(t *testing.T) {
	userID := uuid.New()
	feedID := uuid.New()

	tests := []struct {
		name        string
		setupDB     func(*MockDb)
		params      database.CreateFeedFollowParams
		expectError bool
	}{
		{
			name: "create feed follow successfully",
			setupDB: func(db *MockDb) {
				db.Users["alice"] = database.User{
					ID:        userID,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Name:      "alice",
				}
				db.Feeds["https://example.com/feed"] = database.Rssfeed{
					ID:        feedID,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Name:      "Example Feed",
					Url:       "https://example.com/feed",
					UserID:    userID,
				}
			},
			params: database.CreateFeedFollowParams{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				UserID:    userID,
				FeedID:    feedID,
			},
			expectError: false,
		},
		{
			name:    "create feed follow with non-existing feed",
			setupDB: func(db *MockDb) {},
			params: database.CreateFeedFollowParams{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				UserID:    uuid.New(),
				FeedID:    uuid.New(),
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb := NewMockDb()
			tt.setupDB(mockDb)

			result, err := mockDb.CreateFeedFollow(context.Background(), tt.params)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
				if result.UserID != tt.params.UserID {
					t.Errorf("expected userID %v but got %v", tt.params.UserID, result.UserID)
				}
				if result.FeedID != tt.params.FeedID {
					t.Errorf("expected feedID %v but got %v", tt.params.FeedID, result.FeedID)
				}
			}
		})
	}
}

func TestMockDb_GetFeedFollowsForUser(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name          string
		setupDB       func(*MockDb)
		userID        uuid.UUID
		expectedCount int
	}{
		{
			name:          "get feed follows for user with no follows",
			setupDB:       func(db *MockDb) {},
			userID:        userID,
			expectedCount: 0,
		},
		{
			name: "get feed follows for user with follows",
			setupDB: func(db *MockDb) {
				db.FeedFollows[userID] = []database.GetFeedFollowsForUserRow{
					{
						ID:        uuid.New(),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						UserID:    userID,
						FeedID:    uuid.New(),
						FeedName:  "Feed 1",
						Username:  "alice",
					},
					{
						ID:        uuid.New(),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						UserID:    userID,
						FeedID:    uuid.New(),
						FeedName:  "Feed 2",
						Username:  "alice",
					},
				}
			},
			userID:        userID,
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb := NewMockDb()
			tt.setupDB(mockDb)

			follows, err := mockDb.GetFeedFollowsForUser(context.Background(), tt.userID)

			if err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
			if len(follows) != tt.expectedCount {
				t.Errorf("expected %d follows but got %d", tt.expectedCount, len(follows))
			}
		})
	}
}

func TestMockDb_UnfollowFeed(t *testing.T) {
	userID := uuid.New()
	feedID1 := uuid.New()
	feedID2 := uuid.New()

	tests := []struct {
		name               string
		setupDB            func(*MockDb)
		params             database.UnfollowFeedParams
		expectedFollowsNum int
	}{
		{
			name: "unfollow existing feed",
			setupDB: func(db *MockDb) {
				db.FeedFollows[userID] = []database.GetFeedFollowsForUserRow{
					{
						ID:        uuid.New(),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						UserID:    userID,
						FeedID:    feedID1,
						FeedName:  "Feed 1",
						Username:  "alice",
					},
					{
						ID:        uuid.New(),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						UserID:    userID,
						FeedID:    feedID2,
						FeedName:  "Feed 2",
						Username:  "alice",
					},
				}
			},
			params: database.UnfollowFeedParams{
				UserID: userID,
				FeedID: feedID1,
			},
			expectedFollowsNum: 1,
		},
		{
			name:    "unfollow feed with no existing follows",
			setupDB: func(db *MockDb) {},
			params: database.UnfollowFeedParams{
				UserID: userID,
				FeedID: feedID1,
			},
			expectedFollowsNum: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb := NewMockDb()
			tt.setupDB(mockDb)

			err := mockDb.UnfollowFeed(context.Background(), tt.params)

			if err != nil {
				t.Errorf("expected no error but got: %v", err)
			}

			follows := mockDb.FeedFollows[tt.params.UserID]
			if len(follows) != tt.expectedFollowsNum {
				t.Errorf("expected %d follows remaining but got %d", tt.expectedFollowsNum, len(follows))
			}
		})
	}
}

func TestMockDb_OtherMethods(t *testing.T) {
	mockDb := NewMockDb()

	// Test MarkFeedFetched
	err := mockDb.MarkFeedFetched(context.Background(), uuid.New())
	if err != nil {
		t.Errorf("MarkFeedFetched: expected no error but got: %v", err)
	}

	// Test GetNextFeedToFetch
	_, err = mockDb.GetNextFeedToFetch(context.Background())
	if err != nil {
		t.Errorf("GetNextFeedToFetch: expected no error but got: %v", err)
	}

	// Test CreatePost
	_, err = mockDb.CreatePost(context.Background(), database.CreatePostParams{})
	if err != nil {
		t.Errorf("CreatePost: expected no error but got: %v", err)
	}

	// Test GetPostsForUser
	_, err = mockDb.GetPostsForUser(context.Background(), database.GetPostsForUserParams{})
	if err != nil {
		t.Errorf("GetPostsForUser: expected no error but got: %v", err)
	}

	// Test CreateRSSFeed
	_, err = mockDb.CreateRSSFeed(context.Background(), database.CreateRSSFeedParams{})
	if err != nil {
		t.Errorf("CreateRSSFeed: expected no error but got: %v", err)
	}

	// Test GetFeeds
	_, err = mockDb.GetFeeds(context.Background())
	if err != nil {
		t.Errorf("GetFeeds: expected no error but got: %v", err)
	}
}
