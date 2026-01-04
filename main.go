package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("PORT is not set")
	}

	mainRouter := newRouter()
	server := newServer(":"+port, mainRouter)

	log.Printf("Server is running on port %v :) >> ", port)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
