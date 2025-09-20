package cli

import (
	"context"
	"fmt"
	"log"
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
		log.Fatal("user is not registered in database")
	}

	err = s.Cfg.SetUser(cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("login handler couldn't switch user: %w", err)
	}

	fmt.Println("User login successfull!")
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	_, err := s.Db.GetUser(context.Background(), cmd.Arguments[0])
	if err == nil {
		log.Fatal("user already exists")
	}

	_, err = s.Db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.Arguments[0]})
	if err != nil {
		log.Fatalf("couldn't create new user: %v", err)
	}

	err = s.Cfg.SetUser(cmd.Arguments[0])
	if err != nil {
		log.Fatalf("couldn't set new user in config %v", err)
	}
	fmt.Println("User successfully created!")

	return nil
}
