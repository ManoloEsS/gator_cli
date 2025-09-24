package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/ManoloEsS/gator_cli/cli"
	"github.com/ManoloEsS/gator_cli/internal/config"
	"github.com/ManoloEsS/gator_cli/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	//parse command line arguments
	if len(os.Args) < 2 {
		log.Fatal("argument needed, usage: gator <argument>\n")

	}
	name, args := os.Args[1], os.Args[2:]

	cmd := cli.Command{
		Name:      name,
		Arguments: args,
	}

	//read config file
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	//open database using url in config
	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatalf("database could not be opened: %v", err)
	}
	defer db.Close()

	//ping database to check if it's reachable
	err = db.Ping()
	if err != nil {
		log.Fatalf("database could not be reached: %v", err)
	}

	//initialize State and Commands from cli
	//add database and its function=>queries to program state
	dbQueries := database.New(db)

	programState := cli.NewState(dbQueries, &cfg)

	cmds := cli.Commands{
		CommandMap: make(map[string]func(*cli.State, cli.Command) error),
	}

	//register each command in commands
	cmds.Register("login", cli.HandlerLogin)
	cmds.Register("register", cli.HandlerRegister)
	cmds.Register("reset", cli.HandlerReset)
	cmds.Register("users", cli.HandlerListUsers)
	cmds.Register("agg", cli.HandlerAgg)
	cmds.Register("addfeed", cli.HandlerAddFeed)
	cmds.Register("feeds", cli.HandlerListFeeds)

	//run command from parsed command line arguments
	err = cmds.Run(programState, cmd)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
