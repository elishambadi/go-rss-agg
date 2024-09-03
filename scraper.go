package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/elishambadi/go-rss-agg/internal/database"
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

	for _, item := range rssFeed.Channel.Item {
		log.Println("Found post", item.Title, " on feed", feed.Name)
	}

	log.Printf("Feed %v collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}
