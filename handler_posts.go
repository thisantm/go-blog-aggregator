package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/thisantm/go-blog-aggregator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2

	if len(cmd.args) > 1 {
		return fmt.Errorf("the browse command accepts 1 optional argument ex.: gator browse <limit>")
	}

	if len(cmd.args) == 1 {
		var err error
		limit, err = strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("couldn't parse the limit %v: %v", cmd.args[0], err)
		}
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		ID:    user.ID,
		Limit: int32(limit),
	})
	if err != nil {
		fmt.Printf("failed to get posts user follows: %v", err)
	}

	postsJson, _ := json.MarshalIndent(posts, "", "  ")
	fmt.Println("feed follow: " + string(postsJson))

	return nil
}
