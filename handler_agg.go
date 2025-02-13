package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/thisantm/go-blog-aggregator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("the agg command expects 1 argument ex.: gator agg <time_between_reqs> which is a time.duration ex.: 1s, 1m, 1h")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("invalid time duration: %v", err)
	}

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			return fmt.Errorf("error while scraping: %v", err)
		}
	}
}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %v", err)
	}

	s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		UpdatedAt: time.Now().UTC(),
		LastFetchedAt: sql.NullTime{
			Time:  time.Now().UTC(),
			Valid: true,
		},
		ID: feed.ID,
	})

	fmt.Printf("fetching feed: %s at time: %v\n", feed.Url, time.Now().Format(time.RFC822Z))
	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("failed to fetch rss feed: %v", err)
	}

	rssFeedJson, _ := json.MarshalIndent(rssFeed, "", "  ")
	fmt.Println(string(rssFeedJson))

	return nil
}
