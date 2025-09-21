package cli

import (
	"context"
	"github.com/ManoloEsS/gator_cli/internal/config"
	"github.com/ManoloEsS/gator_cli/internal/database"
)

// DBInterface defines the database operations needed by handlers
type DBInterface interface {
	GetUser(ctx context.Context, name string) (database.User, error)
	CreateUser(ctx context.Context, params database.CreateUserParams) (database.User, error)
	ResetUsers(ctx context.Context) error
}

// ConfigInterface defines the config operations needed by handlers
type ConfigInterface interface {
	SetUser(name string) error
}

type State struct {
	Db  DBInterface
	Cfg ConfigInterface
}

// NewState creates a new State with concrete types
func NewState(db *database.Queries, cfg *config.Config) *State {
	return &State{
		Db:  db,
		Cfg: cfg,
	}
}
