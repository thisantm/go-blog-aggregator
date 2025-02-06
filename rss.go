package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {
	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, feedUrl, http.NoBody)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("failed to write request: %v", err)
	}
	request.Header.Set("User-Agent", "gator")

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("request failed: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("failed to read body: %v", err)
	}

	var rssFeed RSSFeed
	err = xml.Unmarshal(body, &rssFeed)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("failed to unmarshal xml: %v", err)
	}

	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	for i, v := range rssFeed.Channel.Item {
		rssFeed.Channel.Item[i].Title = html.UnescapeString(v.Title)
		rssFeed.Channel.Item[i].Description = html.UnescapeString(v.Description)
	}

	return &rssFeed, nil
}
