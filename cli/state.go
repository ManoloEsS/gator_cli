package cli

import (
	"context"

	"github.com/ManoloEsS/gator_cli/internal/config"
	"github.com/ManoloEsS/gator_cli/internal/database"
)

//Db and Config interfaces are used to be able to test
//creating mocks that fit the interface schema
//while keeping functionality in the program

// DBInterface defines the database operations needed for Database Interface
type DBInterface interface {
	GetUser(ctx context.Context, name string) (database.User, error)
	CreateUser(ctx context.Context, params database.CreateUserParams) (database.User, error)
	ResetUsers(ctx context.Context) error
	GetUsers(ctx context.Context) ([]database.User, error)
	CreateRSSFeed(ctx context.Context, arg database.CreateRSSFeedParams) (database.Rssfeed, error)
	GetFeeds(ctx context.Context) ([]database.GetFeedsRow, error)
	GetFeedByUrl(ctx context.Context, url string) (database.Rssfeed, error)
}

// ConfigInterface defines the config operations needed by Config Interface
type ConfigInterface interface {
	SetUser(name string) error
	GetCurrentUser() string
}

// State struct that stores the Database and Config interfaces
// throughout the program
type State struct {
	Db  DBInterface
	Cfg ConfigInterface
}

// NewState creates a new State with the interfaces
func NewState(db *database.Queries, cfg *config.Config) *State {
	return &State{
		Db:  db,
		Cfg: cfg,
	}
}
