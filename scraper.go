package main

import (
	"context"
	"log"
	"sync"
	"time"

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
		log.Println("Processing item:", item.Title)
	}
	log.Println("Finished scraping feed:", feed.Name)
}
