package main

import (
	"context"
	"encoding/json"
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

	bytes, _ := json.MarshalIndent(user, "", "  ")
	fmt.Println(string(bytes))

	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("the reset command expects no argument ex.: gator reset")
	}

	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to delete users from database: %v", err)
	}

	return nil
}

func handlerUsers(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("the users command expects no arguments ex.: gator users")
	}
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get users from database: %v", err)
	}

	for _, u := range users {
		if u.Name == s.config.CurrentUserName {
			fmt.Printf("* %s (current)\n", u.Name)
			continue
		}
		fmt.Printf("* %s\n", u.Name)
	}

	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("the agg command expects no arguments ex.: gator agg")
	}

	rssFeed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %v", rssFeed)
	}

	rssFeedJson, _ := json.MarshalIndent(rssFeed, "", "  ")
	fmt.Println(string(rssFeedJson))

	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("the addfeed command expects 2 arguments ex.: gator addfeed <name> <url>")
	}

	name := cmd.args[0]
	url := cmd.args[1]

	user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("failed to get user from database: %v", err)
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create feed: %v", err)
	}

	feedJson, _ := json.MarshalIndent(feed, "", "  ")
	fmt.Println(string(feedJson))

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("the feeds command expects no arguments ex.: gator feeds")
	}
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get feeds from database: %v", err)
	}

	for _, f := range feeds {
		fmt.Printf("Feed Name: %s\n", f.FeedName)
		fmt.Printf("Url: %s\n", f.Url)
		fmt.Printf("User Name: %s\n\n", f.UserName)
	}

	return nil
}
