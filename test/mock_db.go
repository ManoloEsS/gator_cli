package test

import (
	"context"
	"errors"

	"github.com/ManoloEsS/gator_cli/internal/database"
	"github.com/google/uuid"
)

// Mock implementations for testing
type MockDb struct {
	Users        map[string]database.User
	Feeds        map[string]database.Rssfeed
	FeedFollows  map[uuid.UUID][]database.GetFeedFollowsForUserRow
	CreateError  error
	ResetError   error
	GetFeedError error
}

func NewMockDb() *MockDb {
	return &MockDb{
		Users:       make(map[string]database.User),
		Feeds:       make(map[string]database.Rssfeed),
		FeedFollows: make(map[uuid.UUID][]database.GetFeedFollowsForUserRow),
	}
}

func (m *MockDb) GetUser(ctx context.Context, name string) (database.User, error) {
	user, exists := m.Users[name]
	if !exists {
		return database.User{}, errors.New("sql: no rows in result set")
	}
	return user, nil
}

// TODO:Add tests for GetUsers method
func (m *MockDb) GetUsers(ctx context.Context) ([]database.User, error) {
	userList := []database.User{}
	for _, user := range m.Users {
		userList = append(userList, user)
	}

	return userList, nil
}

func (m *MockDb) CreateUser(ctx context.Context, params database.CreateUserParams) (database.User, error) {
	if m.CreateError != nil {
		return database.User{}, m.CreateError
	}
	if _, exists := m.Users[params.Name]; exists {
		return database.User{}, errors.New("user already exists")
	}
	user := database.User{
		ID:        params.ID,
		CreatedAt: params.CreatedAt,
		UpdatedAt: params.UpdatedAt,
		Name:      params.Name,
	}
	m.Users[params.Name] = user
	return user, nil
}

func (m *MockDb) ResetUsers(ctx context.Context) error {
	if m.ResetError != nil {
		return m.ResetError
	}
	m.Users = make(map[string]database.User)
	return nil
}

// TODO:finish function
func (m *MockDb) CreateRSSFeed(ctx context.Context, arg database.CreateRSSFeedParams) (database.CreateRSSFeedRow, error) {
	return database.CreateRSSFeedRow{}, nil
}

// TODO:finish test function
func (m *MockDb) GetFeeds(ctx context.Context) ([]database.GetFeedsRow, error) {
	return []database.GetFeedsRow{}, nil
}

func (m *MockDb) GetFeedByUrl(ctx context.Context, url string) (database.Rssfeed, error) {
	if m.GetFeedError != nil {
		return database.Rssfeed{}, m.GetFeedError
	}
	feed, exists := m.Feeds[url]
	if !exists {
		return database.Rssfeed{}, errors.New("feed not found")
	}
	return feed, nil
}

func (m *MockDb) CreateFeedFollow(ctx context.Context, args database.CreateFeedFollowParams) (database.CreateFeedFollowRow, error) {
	// Get feed and user to populate the result
	feed, feedExists := m.Feeds[getFeedUrlByID(m.Feeds, args.FeedID)]
	if !feedExists {
		return database.CreateFeedFollowRow{}, errors.New("feed not found")
	}
	
	user, userExists := getUserByID(m.Users, args.UserID)
	if !userExists {
		return database.CreateFeedFollowRow{}, errors.New("user not found")
	}

	// Add to feed follows
	followRow := database.GetFeedFollowsForUserRow{
		ID:        args.ID,
		CreatedAt: args.CreatedAt,
		UpdatedAt: args.UpdatedAt,
		UserID:    args.UserID,
		FeedID:    args.FeedID,
		FeedName:  feed.Name,
		Username:  user.Name,
	}
	m.FeedFollows[args.UserID] = append(m.FeedFollows[args.UserID], followRow)

	return database.CreateFeedFollowRow{
		ID:        args.ID,
		CreatedAt: args.CreatedAt,
		UpdatedAt: args.UpdatedAt,
		UserID:    args.UserID,
		FeedID:    args.FeedID,
		FeedName:  feed.Name,
		UserName:  user.Name,
	}, nil
}

func (m *MockDb) GetFeedFollowsForUser(ctx context.Context, userID uuid.UUID) ([]database.GetFeedFollowsForUserRow, error) {
	follows, exists := m.FeedFollows[userID]
	if !exists {
		return []database.GetFeedFollowsForUserRow{}, nil
	}
	return follows, nil
}

func (m *MockDb) UnfollowFeed(ctx context.Context, arg database.UnfollowFeedParams) error {
	follows, exists := m.FeedFollows[arg.UserID]
	if !exists {
		return nil
	}

	// Remove the feed follow
	newFollows := []database.GetFeedFollowsForUserRow{}
	for _, follow := range follows {
		if follow.FeedID != arg.FeedID {
			newFollows = append(newFollows, follow)
		}
	}
	m.FeedFollows[arg.UserID] = newFollows
	return nil
}

// Helper functions
func getFeedUrlByID(feeds map[string]database.Rssfeed, feedID uuid.UUID) string {
	for url, feed := range feeds {
		if feed.ID == feedID {
			return url
		}
	}
	return ""
}

func getUserByID(users map[string]database.User, userID uuid.UUID) (database.User, bool) {
	for _, user := range users {
		if user.ID == userID {
			return user, true
		}
	}
	return database.User{}, false
}

func (m *MockDb) MarkFeedFetched(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (m *MockDb) GetNextFeedToFetch(ctx context.Context) (database.Rssfeed, error) {
	return database.Rssfeed{}, nil
}

func (m *MockDb) CreatePost(ctx context.Context, arg database.CreatePostParams) (database.Post, error) {
	return database.Post{}, nil
}

func (m *MockDb) GetPostsForUser(ctx context.Context, arg database.GetPostsForUserParams) ([]database.GetPostsForUserRow, error) {
	return []database.GetPostsForUserRow{}, nil
}
