package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/elishambadi/go-rss-agg/internal/database"
	"github.com/google/uuid"
)

// Scraper to run on multiple different goroutines every few minutes or so
func startScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Scraping on %v goroutines every %s duration.", concurrency, timeBetweenRequest)

	ticker := time.NewTicker(timeBetweenRequest)
	// Ticker sends a value(a tick like a ping) across the channel every timeBetweenRequest

	// Recevie the value from ticker
	for ; ; <-ticker.C {
		feeds, err := db.GetNextsFeedToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Printf("error fetching feeds: %s", err)
			continue
		}

		wg := &sync.WaitGroup{}
		// You can spawn a goroutine in the waitgroups

		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(db, wg, feed)
		}

		wg.Wait() // Waits for all goroutines to finish
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking feed as fetched: ", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching Feed: ", err)
		return
	}

	type RSSItem struct {
		Title       string `xml:"title"`
		Link        string `xml:"link"`
		Description string `xml:"description"`
		PubDate     string `xml:"pubDate"`
	}

	for _, item := range rssFeed.Channel.Item {
		log.Println("Found post", item.Title, " on feed", feed.Name)

		// Handling the NULL value safely
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		// Handling a time parsing
		pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("couldn't parse date %v with err %v", item.PubDate, err)
			continue
		}

		// Store the items in DB
		_, err1 := db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: description,
			PublishedAt: pubAt,
			Url:         item.Link,
			FeedID:      feed.ID,
		})

		if err1 != nil {
			if strings.Contains(err1.Error(), "duplicate key value") {
				continue
			}
			log.Printf("Error while creating post: %v", err1)
		}
	}

	log.Printf("Feed %v collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}
