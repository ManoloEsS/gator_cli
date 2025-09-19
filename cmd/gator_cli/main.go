package main

import (
	"fmt"
	"os"

	"github.com/ManoloEsS/gator_cli/cli"
	"github.com/ManoloEsS/gator_cli/internal/config"
)

func main() {
	//read config file
	cfg, err := config.Read()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	//initialize State and Commands from cli
	programState := cli.State{}
	cmds := cli.Commands{
		CommandMap: make(map[string]func(*cli.State, cli.Command) error),
	}

	//save read config to current state
	programState.Cfg = &cfg
	//register login command in commands
	cmds.Register("login", cli.HandlerLogin)

	if len(os.Args) < 2 {
		fmt.Println("argument needed, usage: gator <argument>")
		os.Exit(1)
	}
	name, args := os.Args[1], os.Args[2:]

	cmd := cli.Command{
		Name:      name,
		Arguments: args,
	}

	err = cmds.Run(&programState, cmd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
