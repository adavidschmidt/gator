package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/adavidschmidt/blogaggregator/internal/database"
	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("agg <time_between_reqs> (eg '1m', '1s', '1h')")
	}
	d, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Collecting feeds every %v\n", d)
	ticker := time.NewTicker(d)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) error {
	now := time.Now().UTC()
	fetchedFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		UpdatedAt: now,
		ID:        fetchedFeed.ID})
	if err != nil {
		return err
	}
	feed, err := fetchFeed(context.Background(), fetchedFeed.Url)
	if err != nil {
		return err
	}
	for _, item := range feed.Channel.Item {
		t, err := parsePubDate(item.PubDate)
		if err != nil {
			continue
		}
		err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   now,
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: t,
			FeedID:      fetchedFeed.ID,
		})
		if err != nil {
			continue
		}
	}

	return nil
}

func parsePubDate(s string) (time.Time, error) {
	layouts := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		time.RFC3339,
		"2006-01-02 15:04:05",
		"Mon, 02 Jan 2006 15:04:05 -0700",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unrecognized date format: %s", s)
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int64
	limit = 2
	var err error
	if len(cmd.Args) > 0 {
		limit, err = strconv.ParseInt(cmd.Args[0], 10, 32)
		if err != nil {
			return fmt.Errorf("browse <number of posts to browse>")
		}
	}
	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	for _, post := range posts {
		fmt.Printf("%s from %s\n", post.PublishedAt.Format("Mon Jan 2"), post.FeedName)
		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("    %v\n", post.Description)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("=====================================")
	}
	return nil
}
