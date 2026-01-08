package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alepaez-dev/rss_aggregator/internal/database"
	"github.com/google/uuid"
)

type resp struct {
	Status string `json:"status"` // w this tag -> instead of "Status" it will be json format "status"
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	err := respondWithJSON(w, http.StatusOK, resp{Status: "OK server good 🧃 :)"})
	if err != nil {
		log.Printf("Error writing JSON response: %v", err)
	}
}

func handlerErr(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 400, "Something went wrong")
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	// Go treats http bodies as STREAMS instead of buffers -> we need to handle it. (fast for RAM 🐏)

	type parameters struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body) // decoder reads directly from the http body stream.
	err := decoder.Decode(&params)     // maps JSON to struct(params) -- a pointer is needed to achieve this.
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FirstName: params.FirstName,
		LastName:  params.LastName,
	})

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Couldn't create user: %v", err))
		return
	}

	respondWithJSON(
		w,
		http.StatusCreated,
		resp{Status: fmt.Sprintf("User creates successfullu %v 🧃 :)", user.FirstName)},
	)
}
