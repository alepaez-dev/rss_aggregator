package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

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

func (cfg *ApiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
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
		FirstName: params.FirstName,
		LastName:  params.LastName,
	})

	if err != nil {
		log.Printf("Couldn't create user %v", err)
		respondWithError(w, http.StatusBadRequest, "Couldn't create user")
		return
	}

	respondWithJSON(
		w,
		http.StatusCreated,
		databaseUserToUser(user),
	)
}

func (cfg *ApiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}

// Get the posts of the feed that the user follows
func (cfg *ApiConfig) handlerGetPostsForUser(w http.ResponseWriter, r *http.Request, user database.User) {
	limit := 2
	if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
		parsedLimit, err := strconv.Atoi(limitParam)
		if err != nil || parsedLimit <= 0 {
			respondWithError(w, http.StatusBadRequest, "Invalid limit")
			return
		}
		limit = parsedLimit
	}

	posts, err := cfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't get post for user")
		return
	}

	respondWithJSON(w, http.StatusOK, databasePostsToPosts(posts))
}
