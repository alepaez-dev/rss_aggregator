package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func NewRouter(cfg *ApiConfig) http.Handler {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000", "https://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
		},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	// Base
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)

	// Users
	v1Router.Post("/users", cfg.handlerCreateUser)
	v1Router.Get("/users", cfg.middlewareAuth(cfg.handlerGetUser))

	// Feeds
	v1Router.Post("/feeds", cfg.middlewareAuth(cfg.handlerCreateFeed))
	v1Router.Get("/feeds", cfg.handlerGetFeeds)

	// Feeds Follows
	v1Router.Post("/feed_follows", cfg.middlewareAuth(cfg.handlerCreateFeedFollow))
	v1Router.Get("/feed_follows", cfg.middlewareAuth(cfg.handlerGetFeedFollows))
	v1Router.Delete("/feed_follows/{feedFollowID}", cfg.middlewareAuth(cfg.handlerDeleteFeedFollow))

	// V1
	r.Mount("/v1", v1Router)

	return r
}
