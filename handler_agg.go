package main

import (
	"context"
	"encoding/json"
	"fmt"
)

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
