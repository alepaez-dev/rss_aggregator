package tasks

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/alepaez-dev/rss_aggregator/internal/database"
	"github.com/alepaez-dev/rss_aggregator/internal/feeds"
)

func scrapeFeed(ctx context.Context, db *database.Queries, feed database.Feed) {
	fmt.Println("Scraping feed:", feed.ID)

	_, err := db.MarkFeedAsFetched(ctx, feed.ID)
	if err != nil {
		log.Printf("Error marking feed as fetched: %v", err)
		return
	}

	rssFeed, err := feeds.UrlToFeed(ctx, feed.Url)
	if err != nil {
		log.Printf("Error fetching feed URL %s: %v", feed.Url, err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		log.Println("Found post", item.Title)
	}

}

func worker(
	ctx context.Context,
	jobs <-chan database.Feed,
	wg *sync.WaitGroup,
	db *database.Queries,
) {

	defer wg.Done() // worker finished
	for {
		select {
		case <-ctx.Done():
			return
		case feed, ok := <-jobs: // we received a feed 🙏
			if !ok { // safe check → is channel closed?
				return
			}
			feedCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
			scrapeFeed(feedCtx, db, feed) // sync
			cancel()                      // free resources, scrapFeed is sync this means it's done
		}
	}
}

func StartScraping(ctx context.Context, db *database.Queries, concurrency int, interval time.Duration) {
	ticker := time.NewTicker(interval)

	// cleanup (2nd)
	defer ticker.Stop()

	jobs := make(chan database.Feed)

	var wg sync.WaitGroup
	wg.Add(concurrency) // we will wait for N workers to finish

	for i := 0; i < concurrency; i++ {
		go worker(ctx, jobs, &wg, db)
	}

	// cleanup (1st)
	defer func() {
		close(jobs) // current jobs keep going but no future jobs
		wg.Wait()   // WAIT for all workers to finish
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			feeds, err := db.GetNextFeedsToFetch(ctx, int32(concurrency))
			if err != nil {
				log.Printf("error fetching feeds: %v", err)
				continue
			}
			for _, f := range feeds {
				select {
				case <-ctx.Done(): // if ctx is done, stop immediately, in case of deadlock
					return
				case jobs <- f: // send job when worker is ready to receive
				}
			}
		}
	}
}
