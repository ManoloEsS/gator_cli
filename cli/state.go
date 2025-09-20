package cli

import (
	"github.com/ManoloEsS/gator_cli/internal/config"
	"github.com/ManoloEsS/gator_cli/internal/database"
)

type State struct {
	Db  *database.Queries
	Cfg *config.Config
}
