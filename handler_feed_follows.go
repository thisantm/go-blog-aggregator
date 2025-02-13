package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/thisantm/go-blog-aggregator/internal/database"
)

func handlerFollow(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("the follow command expects 1 argument ex.: gator follow <url>")
	}

	url := cmd.args[0]

	user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("failed to get user from database: %v", err)
	}

	feed, err := s.db.GetFeedsByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("failed to get feed from database: %v", err)
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create feed_follow entry: %v", err)
	}

	feedFollowJson, _ := json.MarshalIndent(feedFollow, "", "  ")
	fmt.Println(string(feedFollowJson))

	return nil
}

func handlerFollowing(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("the follow command expects 0 argument ex.: gator following")
	}

	user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("failed to get user from database: %v", err)
	}

	feedFollow, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("failed to get feed_follows from database: %v", err)
	}

	feedFollowJson, _ := json.MarshalIndent(feedFollow, "", "  ")
	fmt.Println(string(feedFollowJson))

	return nil
}
