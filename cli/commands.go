package cli

import "errors"

type Command struct {
	Name      string
	Arguments []string
}

type Commands struct {
	CommandMap map[string]func(*State, Command) error
}

func (c *Commands) Run(s *State, cmd Command) error {
	function, ok := c.CommandMap[cmd.Name]
	if !ok {
		return errors.New("command not found")
	}

	return function(s, cmd)
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.CommandMap[name] = f

}
