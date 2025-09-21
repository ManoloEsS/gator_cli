package cli

import (
	"context"
	"fmt"
	"log"
)

func HandlerReset(s *State, cmd Command) error {
	err := s.Db.ResetUsers(context.Background())
	if err != nil {
		log.Fatalf("couldn't reset users table: %v", err)
	}

	fmt.Println("users table has been reset...")
	return nil
}
