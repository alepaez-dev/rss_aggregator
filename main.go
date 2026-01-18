package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alepaez-dev/rss_aggregator/internal/api"
	"github.com/alepaez-dev/rss_aggregator/internal/database"
	"github.com/alepaez-dev/rss_aggregator/internal/tasks"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // postgres driver
)

/*
EXPLANATION OF THE FLOW:
1. We create root context in main
2. We subscribe to OS shutdown signals (Ctrl-C) and listen for those signals to do a graceful shutdown in sigCh channel which block (<-sigCh on line 85) needs to be after the goroutines so we don;t block them.
3. We create DB connection, router, etc
4. We do a scrapeDone channel that will block main program on line 92 (at the end of main). Unless we send a signal here main will never shutdown (unless server never starts, log.Fatal will shutdown everything, is fine).
5. We start the StartScraping in a goroutine(async) with the root context. When the StartScraping finishes synchronously it closes the scrapeDone channel unblocking main program on line 93.
6. We start the HTTP server in another goroutine(async).
7. Once ctrl-c is done <-sigCh is unblocked on line 85 we continue with line 86 execution which is cancel(), it will cancel the root context which will propagate to StartScraping and all of it's child workers.
8. StartScraping exits only after all workers finish. Everything is done gracefully there.
9. The main program will wait for StartScraping to finish on line 93 <-scrapeDone
10. Before we reach line 93 (StartScraping is currently finishing here) we shut down server gracefully with a timeout context of 10 seconds.
11. Once the server finishes or timeout is reached we go to line 93 where main waits for StartScraping to finish if it hasn't already.
*/
func main() {
	// Root context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Listen for OS shutdown signals (ctrl-c, etc)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	godotenv.Load(".env")
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("PORT is not set in environment")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not set in environment")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database", err)
	}
	defer conn.Close()

	queries := database.New(conn)
	cfg := api.ApiConfig{
		DB: queries,
	}

	scrapeDone := make(chan struct{})

	// async
	go func() {
		tasks.StartScraping(ctx, queries, 5, 10*time.Second)
		close(scrapeDone)
	}()

	mainRouter := api.NewRouter(&cfg)
	server := newServer(":"+port, mainRouter)
	// async
	go func() {
		// run server
		log.Printf("Server is running on port %v :) >> ", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err) // TODO: this can be handled better.
		}
	}()

	<-sigCh  // wait for ctrl-c or killl (main is blocker here we can't finish program)
	cancel() // cancel root context, it will propagate to all children in StartScraping

	// need new context to use timeout bcs root context is cancelled already at this line
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	_ = server.Shutdown(shutdownCtx) // server has 10 seconds to shutdown gracefully or bye

	<-scrapeDone // wait for scraping to finish, and StartScraping is waiting for it's workers to finish. Even a Ctrl-C or SIGTERM/SIGINT shutdown wll handle everything gracefully. (not a kill -9 tho)

}
