package tasks

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/alepaez-dev/rss_aggregator/internal/database"
)

func worker(
	ctx context.Context,
	jobs <-chan database.Feed,
) {
	for {
		select {
		case <-ctx.Done():
			return
		case feed, ok := <-jobs:
			if !ok { // safe
				return
			}
			// TODO: implement the actual scraping logic here
			fmt.Println("Processing feed:", feed.ID)
		}
	}
}

func StartScraping(ctx context.Context, db *database.Queries, concurrency int, interval time.Duration) {
	ticker := time.NewTicker(interval)

	defer ticker.Stop()

	jobs := make(chan database.Feed)

	for i := 0; i < concurrency; i++ {
		go worker(ctx, jobs)
	}

	for {
		select {
		case <-ctx.Done():
			close(jobs)
			return
		case <-ticker.C:
			feeds, err := db.GetNextFeedsToFetch(ctx, int32(concurrency))
			if err != nil {
				log.Printf("error fetching feeds: %v", err)
				continue
			}
			for _, f := range feeds {
				jobs <- f
			}
		}
	}
}
