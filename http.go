package main

import (
	"context"
	"net/http"
	"io"
	"encoding/xml"
	"log"
	"html"
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

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := &http.Client{}
	var storage RSSFeed
	requestContent, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		log.Fatalf("request creation failed: %w", err)
	}
	requestContent.Header.Set("User-Agent", "gator")
	responseContent, err := client.Do(requestContent)
	if err != nil {
		log.Fatalf("response return failed: %w", err)
	}
	feedBody, err := io.ReadAll(responseContent.Body)
	if err != nil {
		log.Fatalf("response body read failed: %w", err)
	}
	
	err = xml.Unmarshal(feedBody, &storage)
	if err != nil {
		log.Fatalf("something went wrong unmarshalling xml: %w", err)
	}
	err = unEscapeChars(&storage)
	if err != nil {
		log.Fatalf("unescape didn't work correctly: %w", err)
	}
	return &storage, nil
}

func unEscapeChars(feed *RSSFeed) error {
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i, instance := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(instance.Title)
		feed.Channel.Item[i].Description = html.UnescapeString(instance.Description)
	}
	return nil
}