package api

import "github.com/alepaez-dev/rss_aggregator/internal/database"

type ApiConfig struct {
	DB *database.Queries // depends on database.Queries struct (TODO: change later to behavioral pattern)
}
