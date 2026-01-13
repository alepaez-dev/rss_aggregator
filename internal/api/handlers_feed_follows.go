package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/alepaez-dev/rss_aggregator/internal/database"
	"github.com/alepaez-dev/rss_aggregator/internal/dberr"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:     uuid.New(),
		UserID: user.ID,
		FeedID: params.FeedID,
	})

	if err != nil {
		if dberr.IsUniqueViolation(err) {
			respondWithError(w, http.StatusConflict, "You already follow this feed")
			return
		}

		// generic error
		log.Printf("Error creating feed follow: %v", err)
		respondWithError(w, http.StatusBadRequest, "Couldn't create the follow of the feed")
		return
	}

	respondWithJSON(
		w,
		http.StatusCreated,
		databaseFeedFollowToFeedFollow(feedFollow),
	)
}

func (cfg *ApiConfig) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {

	feedFollows, err := cfg.DB.GetFeedFollows(r.Context(), user.ID)

	if err != nil {
		// generic error
		log.Printf("Error get feed follows for user_id %v: error=%v", user.ID, err)
		respondWithError(w, http.StatusBadRequest, "Couldn't get feed follows")
		return
	}

	respondWithJSON(
		w,
		http.StatusOK,
		databaseFeedFollowsToFeedFollows(feedFollows),
	)
}

func (cfg *ApiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIDStr := chi.URLParam(r, "feedFollowID")

	feedFollowID, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		respondWithError(w, 400, "Invalid feed follow ID")
	}

	rows, err := cfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowID,
		UserID: user.ID,
	})

	if err != nil {
		// generic error
		log.Printf("Error deleting feed follow feed_id %v: error=%v", feedFollowID, err)
		respondWithError(w, http.StatusBadRequest, "Couldn't delete feed follow")
		return
	}

	if rows == 0 {
		respondWithError(w, http.StatusNotFound, "Not found")
		return
	}

	respondWithJSON(
		w,
		http.StatusOK,
		struct{}{},
	)
}
