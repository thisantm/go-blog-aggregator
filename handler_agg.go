package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
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

	for _, item := range rssFeed.Channel.Item {
		validTime, err := parseMultipleLayoutTime(item.PubDate)
		if err != nil {
			return fmt.Errorf("failed to parse pubData: %v", err)
		}
		s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Title:     item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  len(item.Description) > 0,
			},
			PublishedAt: validTime,
			FeedID:      feed.ID,
		})
	}

	return nil
}

func parseMultipleLayoutTime(timeString string) (time.Time, error) {
	layouts := []string{time.RFC1123Z}

	for _, layout := range layouts {
		validTime, err := time.Parse(layout, timeString)
		if err == nil {
			return validTime, nil
		}
	}
	return time.Time{}, fmt.Errorf("no valid time layout found for: %s", timeString)
}
