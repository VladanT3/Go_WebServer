package main

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/VladanT3/Go_WebServer/internal/database"
	"github.com/google/uuid"
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
        description := sql.NullString{}
        if item.Description != "" {
            description.String = item.Description
            description.Valid = true
        }
        pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
        if err != nil {
            log.Printf("Couldn't convert date: %v\n", err)
            continue
        }

        _, err = db.CreatePost(context.Background(), database.CreatePostParams{
            ID: uuid.New(),
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
            Title: item.Title,
            Description: description,
            PublishedAt: pubAt,
            Url: item.Link,
            FeedID: feed.ID,
        })
        if err != nil {
            log.Printf("Error creating new post: %v\n", err)
        }
    }
    log.Printf("Feed %s collected, %v posts found\n", feed.Name, len(rssFeed.Channel.Item))
}
