package test

import (
	"context"
	"errors"

	"github.com/ManoloEsS/gator_cli/internal/database"
	"github.com/google/uuid"
)

// Mock implementations for testing
type MockDb struct {
	Users       map[string]database.User
	CreateError error
	ResetError  error
}

func NewMockDb() *MockDb {
	return &MockDb{
		Users: make(map[string]database.User),
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
func (m *MockDb) CreateRSSFeed(ctx context.Context, arg database.CreateRSSFeedParams) (database.Rssfeed, error) {
	return database.Rssfeed{}, nil
}

// TODO:finish test function
func (m *MockDb) GetFeeds(ctx context.Context) ([]database.GetFeedsRow, error) {
	return []database.GetFeedsRow{}, nil
}

func (m *MockDb) GetFeedByUrl(ctx context.Context, url string) (database.Rssfeed, error) {
	return database.Rssfeed{}, nil
}

func (m *MockDb) CreateFeedFollow(ctx context.Context, args database.CreateFeedFollowParams) (database.CreateFeedFollowRow, error) {
	return database.CreateFeedFollowRow{}, nil
}

func (m *MockDb) GetFeedFollowsForUser(ctx context.Context, userID uuid.UUID) ([]database.GetFeedFollowsForUserRow, error) {
	return []database.GetFeedFollowsForUserRow{}, nil
}
