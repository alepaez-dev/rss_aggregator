-- name: CreateFeedFollow :one
INSERT INTO feed_follows (id, user_id, feed_id)
VALUES ($1, $2, $3)
RETURNING *;


-- name: GetFeedFollows :many
SELECT * from feed_follows where user_id = $1;

-- name: DeleteFeedFollow :execrows
DELETE FROM feed_follows WHERE id = $1 AND user_id = $2;