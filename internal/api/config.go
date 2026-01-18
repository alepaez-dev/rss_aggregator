package api

import (
	"context"

	"github.com/alepaez-dev/rss_aggregator/internal/database"
	"github.com/google/uuid"
)

type DB interface {
	CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error)
	GetUserByAPIKey(ctx context.Context, apiKey string) (database.User, error)
	CreateFeed(ctx context.Context, arg database.CreateFeedParams) (database.Feed, error)
	GetFeeds(ctx context.Context) ([]database.Feed, error)
	CreateFeedFollow(ctx context.Context, arg database.CreateFeedFollowParams) (database.FeedFollow, error)
	GetFeedFollows(ctx context.Context, userID uuid.UUID) ([]database.FeedFollow, error)
	DeleteFeedFollow(ctx context.Context, arg database.DeleteFeedFollowParams) (int64, error)
	GetPostsForUser(ctx context.Context, arg database.GetPostsForUserParams) ([]database.Post, error)
}

type ApiConfig struct {
	DB DB
}
