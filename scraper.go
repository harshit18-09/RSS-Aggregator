package main

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/harshit18-09/RSS-Aggregator/internal/db"
)

func startscraping(dbq *db.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Println("Scraper started")
	ticker := time.NewTicker(timeBetweenRequest)

	for ; ; <-ticker.C {
		feeds, err := dbq.GetNextFeedToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println("Error fetching feeds to scrape:", err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(dbq, wg, feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(dbq *db.Queries, wg *sync.WaitGroup, feed db.Feed) {
	defer wg.Done()

	err := dbq.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking feed as fetched:", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed from URL:", err)
		return
	}

	for _, item := range rssFeed.Channel.Items {
		t, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			t = time.Now()
		}

		_, err = dbq.CreatePost(context.Background(), db.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			FeedID:      uuid.NullUUID{UUID: feed.ID, Valid: true},
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: item.Description != ""},
			PublishedAt: t,
		})
		if err != nil {
			log.Println("Error inserting post into database:", err)
			continue
		}
	}

	log.Println("Finished scraping feed:", feed.Name)
}
