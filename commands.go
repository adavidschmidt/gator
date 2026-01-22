package main

import (
	"errors"
	"context"

	"github.com/adavidschmidt/blogaggregator/internal/database"
)
type command struct {
	Name string
	Args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.handlers[cmd.Name]
	if !ok {
		return errors.New("Uknown Command\n")
	}

	err := handler(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		loggedInUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUser)
		if err != nil {
			return err
		}
		return handler(s, cmd, loggedInUser)
	}
}

