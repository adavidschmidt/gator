package main

import (
	"context"
	"fmt"
	"time"

	"github.com/adavidschmidt/blogaggregator/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("addfeed <feed name> <feed url>")
	}
	name := cmd.Args[0]
	url := cmd.Args[1]
	now := time.Now().UTC()
	

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}
	fmt.Printf("Feed created: \n%+v\n", feed)
	if err = followFeedByURL(s, url, user); err != nil {
		return err
	}
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("Too many arguments provided")
	}
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found")
		return nil
	}

	for _, feed := range feeds {
		username, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return err
		}

		fmt.Printf("* ID: %s\n", feed.ID)
		fmt.Printf("* Created at: %v\n", feed.CreatedAt)
		fmt.Printf("* Updated at: %v\n", feed.UpdatedAt)
		fmt.Printf("* Name: %s\n", feed.Name)
		fmt.Printf("* URL: %s\n", feed.Url)
		fmt.Printf("* User: %s\n", username.Name)
		fmt.Println("\n==============================================")
	}
	return nil
}
