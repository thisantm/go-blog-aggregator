package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/thisantm/go-blog-aggregator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("the login command expects 1 argument ex.: gator login <user>")
	}

	userName := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), userName)
	if err != nil {
		return fmt.Errorf("user does not exist: %v", err)
	}

	err = s.config.SetUser(userName)
	if err != nil {
		return fmt.Errorf("failed to set user in config file: %v", err)
	}

	fmt.Println("user has been set")

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("the login command expects 1 argument ex.: gator register <user>")
	}

	userName := cmd.args[0]

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      userName,
	})
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	err = s.config.SetUser(userName)
	if err != nil {
		return fmt.Errorf("failed to set current user to new user: %v", err)
	}

	fmt.Printf("%+v\n", user)

	return nil
}
