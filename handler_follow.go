package main

import (
	"context"
	"fmt"
	"time"

	"github.com/adavidschmidt/blogaggregator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("follow <url>")
	}
	return followFeedByURL(s, cmd.Args[0], user)
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	if len(follows) == 0 {
		fmt.Println("No followed feeds")
		return nil
	}
	fmt.Println("Following:")
	for _, follow := range follows {
		fmt.Printf("* %s\n", follow.String)
	}
	return nil
}

func followFeedByURL(s *state, url string, user database.User) error {
	ctx := context.Background()
	now := time.Now().UTC()

	feed, err := s.db.GetFeedByUrl(ctx, url)
	if err != nil {
		return err
	}


	follow, err := s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}

	fmt.Println("Follow created for:")
	fmt.Printf("Feed: %s\n", follow.FeedName)
	fmt.Printf("For user: %s\n", follow.UserName)
	return nil
}


func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("unfollow <url>")
	}
	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}

	err =  s.db.DeleteFollow(context.Background(), database.DeleteFollowParams {
		UserID: user.ID, 
		FeedID: feed.ID})
	if err != nil {
		return err
	}
	fmt.Printf("Successfully unfollowed %s", feed.Name)
	return nil
}