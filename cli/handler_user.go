package cli

import "fmt"

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	err := s.Cfg.SetUser(cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("login handler couldn't set username: %w", err)
	}

	fmt.Println("User set successfully!")
	return nil
}
