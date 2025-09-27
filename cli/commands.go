package cli

import (
	"errors"
)

// signature for cli commands
type Command struct {
	Name      string
	Arguments []string
}

// map that stores registered commands
type Commands struct {
	CommandMap map[string]func(*State, Command) error
}

// handler that runs a command using the os.Args
func (c *Commands) Run(s *State, cmd Command) error {
	function, ok := c.CommandMap[cmd.Name]
	if !ok {
		return errors.New("command not found")
	}

	return function(s, cmd)
}

// function that registers usable commands to the Commands map
func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.CommandMap[name] = f

}
