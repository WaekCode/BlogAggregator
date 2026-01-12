package main

import (
	"context"
	"fmt"
	"log"

	// "fmt"
	"time"

	"github.com/WaekCode/BlogAggregator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func scrapeFeeds(s *State) error {
	fmt.Println()
	fmt.Println("Scrapping Posts from feeds")

	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	err2 := s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err2 != nil {
		return err2
	}

	rssfeed, err3 := fetchFeed(context.Background(), feed.Url)
	if err3 != nil {
		return err3
	}

	for _, item := range rssfeed.Channel.Item {
		if item.Description == "" {
			item.Description = ""
		}

		// Try RFC1123 first (e.g., "Wed, 12 Jan 2026 15:04:05 GMT")
		pubDate, err23 := time.Parse(time.RFC1123, item.PubDate)
		if err23 != nil {
			// Fallback to RFC1123Z (e.g., "Wed, 12 Jan 2026 15:04:05 -0700")
			pubDate, err23 = time.Parse(time.RFC1123Z, item.PubDate)
			if err23 != nil {
				// Log the error and use current time as fallback
				log.Printf("Could not convert PubDate '%s' into time: %v. Using time.Now()", item.PubDate, err23)
				pubDate = time.Now()
			}
		}

		// Create data in the posts table
		_, err4 := s.db.CreatePost(context.Background(),
			database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Title:       item.Title,
				Url:         item.Link,
				Description: item.Description,
				PublishedAt: pubDate,
				FeedID:      feed.ID,
			})
		if err4 != nil {
			if pqErr, ok := err4.(*pq.Error); ok {
				if pqErr.Code == "23505" { // unique violation
					// Duplicate URL — just ignore
					continue
				}
			}
			// Other errors — log or handle
			log.Printf("Failed to insert post: %v", err4)
			continue
		}

	}

	return nil

}
