package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/ManoloEsS/gator_cli/internal/database"
	"github.com/google/uuid"
)

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	_, err := s.Db.GetUser(context.Background(), cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("user is not registered in database\n")
	}

	err = s.Cfg.SetUser(cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("login handler couldn't switch user: %w\n", err)
	}

	fmt.Println("User login successfull!")
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("usage: %s <name>\n", cmd.Name)
	}
	_, err := s.Db.GetUser(context.Background(), cmd.Arguments[0])
	if err == nil {
		return fmt.Errorf("user already exists\n")
	}

	_, err = s.Db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.Arguments[0]})
	if err != nil {
		return fmt.Errorf("couldn't create new user: %v\n", err)
	}

	err = s.Cfg.SetUser(cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("couldn't set new user in config %v\n", err)
	}
	fmt.Println("User successfully created!")

	return nil
}
