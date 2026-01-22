package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/adavidschmidt/blogaggregator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return errors.New("Missing username argument")
	}
	username := cmd.Args[0]
	ctx := context.Background()
	_, err := s.db.GetUser(ctx, username)
	if err != nil {
		return err
	}
	if err = s.cfg.SetUser(username); err != nil {
		return fmt.Errorf("could not set the current user: %w", err)
	}
	fmt.Printf("User %s has been set\n", username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return errors.New("Missing Username")
	}
	name := cmd.Args[0]
	ctx := context.Background()
	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
	}
	user, err := s.db.CreateUser(ctx, params)
	if err != nil {
		return fmt.Errorf("couldn't create user: %w", err)
	}
	err = s.cfg.SetUser(name)
	if err != nil {
		return err
	}
	fmt.Printf("created user:\n%+v\n", user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	ctx := context.Background()
	if err := s.db.ResetUsers(ctx); err != nil {
		return err
	}
	fmt.Println("successfully deleted data from users table")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	if len(cmd.Args) > 0 {
		return errors.New("Too Many Arguments")
	}
	ctx := context.Background()
	users, err := s.db.GetUsers(ctx)
	if err != nil {
		return err
	}
	activeUser := s.cfg.CurrentUser
	for _, user := range users {
		if user.Name == activeUser {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}

