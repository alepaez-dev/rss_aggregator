package api

import (
	"fmt"
	"net/http"

	"github.com/alepaez-dev/rss_aggregator/internal/auth"
	"github.com/alepaez-dev/rss_aggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *ApiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, http.StatusForbidden, fmt.Sprintf("You are NOT authorized: %v", err))
			return
		}

		user, err := cfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, http.StatusForbidden, "You are NOT authorized")
			return
		}

		handler(w, r, user)
	}
}
