package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/alepaez-dev/rss_aggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:     uuid.New(),
		Name:   params.Name,
		Url:    params.Url,
		UserID: user.ID,
	})

	if err != nil {
		log.Printf("Couldn't create feed: %v", err)
		respondWithError(w, http.StatusBadRequest, "Couldn't create feed")
		return
	}

	respondWithJSON(
		w,
		http.StatusCreated,
		databaseFeedToFeed(feed),
	)
}

func (cfg *ApiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfg.DB.GetFeeds(r.Context())

	if err != nil {
		log.Printf("Couldn't get feeds: %v", err)
		respondWithError(w, http.StatusBadRequest, "Couldn't get feeds")
		return
	}

	respondWithJSON(
		w,
		http.StatusAccepted,
		databaseFeedsToFeed(feeds),
	)
}
