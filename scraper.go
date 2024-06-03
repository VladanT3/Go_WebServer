package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/VladanT3/Go_WebServer/internal/database"
)

func startScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration){
    log.Printf("Scraping on %v goroutines every %s duration", concurrency, timeBetweenRequest)
    ticker := time.NewTicker(timeBetweenRequest)
    for ; ; <-ticker.C {
        feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
        if err != nil {
            log.Printf("Error fetching feeds: %v\n", err)
            continue
        }
        
        wg := &sync.WaitGroup{}
        for _, feed := range feeds {
            wg.Add(1)

            go scrapeFeed(wg, db, feed)
        }
        wg.Wait()
    }
}

func scrapeFeed(wg *sync.WaitGroup, db *database.Queries, feed database.Feed) {
    defer wg.Done()

    _, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
    if err != nil {
        log.Printf("Error marking feed as fetched: %v\n", err)
    }

    rssFeed, err := urlToFeed(feed.Url)
    if err != nil {
        log.Printf("Error getting feed from url: %v\n", err)
    }

    for _, item := range rssFeed.Channel.Item {
        log.Printf("Fount post: %v", item.Title)
    }
    log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))

}
