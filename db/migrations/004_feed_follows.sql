-- +goose Up

CREATE TABLE feed_follows (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    feed_id UUID NOT NULL,

    
    UNIQUE (user_id, feed_id) -- Prevent duplicate follows + makes queries (user_id + feed_id OR WHERE user_id = $1) faster
);


 -- Index to find all users that follow a feed without scanning the whole table
 -- feed 1 → user A, user B, user C
 -- feed 2 → user A
 -- feed 3 → user D
CREATE INDEX feed_follows_feed_id_idx
ON feed_follows (feed_id);

-- +goose Down

DROP TABLE feed_follows;
