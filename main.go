package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/alepaez-dev/rss_aggregator/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq" // postgres driver
)

type apiConfig struct {
	DB *database.Queries // depends on database.Queries struct (TODO: change later to behavioral pattern)
}

func main() {
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

	queries := database.New(conn)
	cfg := apiConfig{
		DB: queries,
	}

	mainRouter := newRouter(&cfg)
	server := newServer(":"+port, mainRouter)

	log.Printf("Server is running on port %v :) >> ", port)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
