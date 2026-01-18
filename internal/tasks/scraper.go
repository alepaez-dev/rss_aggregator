package tasks

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/alepaez-dev/rss_aggregator/internal/database"
	"github.com/alepaez-dev/rss_aggregator/internal/dberr"
	"github.com/alepaez-dev/rss_aggregator/internal/feeds"
	"github.com/google/uuid"
)

func scrapeFeed(ctx context.Context, db *database.Queries, feed database.Feed) error {
	_, err := db.MarkFeedAsFetched(ctx, feed.ID)
	if err != nil {
		return fmt.Errorf("mark feed as fetched: %w", err)
	}

	rssFeed, err := feeds.UrlToFeed(ctx, feed.Url)
	if err != nil {
		return fmt.Errorf("fetch feed URL %s: %w", feed.Url, err)
	}

	for _, item := range rssFeed.Channel.Item {
		publishedAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			publishedAt, err = time.Parse(time.RFC1123, item.PubDate)
			if err != nil {
				return fmt.Errorf("parse pubDate %q for feed %s: %w", item.PubDate, feed.ID, err)
			}
		}

		_, err = db.CreatePost(ctx, database.CreatePostParams{
			ID:          uuid.New(),
			Title:       item.Title,
			Description: sql.NullString{String: item.Description, Valid: item.Description != ""},
			PublishedAt: publishedAt,
			Url:         item.Link,
			FeedID:      feed.ID,
		})
		if err != nil {
			if dberr.IsUniqueViolation(err) {
				fmt.Println("Post already exists, skipping")
				continue
			}
			return fmt.Errorf("create post for feed %s: %w", feed.ID, err)
		}
	}

	return nil
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
			time.Sleep(5 * time.Second)
			feedCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
			if err := scrapeFeed(feedCtx, db, feed); err != nil {
				log.Printf("Error scraping feed %s: %v", feed.ID, err)
			}
			cancel() // free resources, scrapFeed is sync this means it's done
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
