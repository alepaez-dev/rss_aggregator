-- name: CreatePost :one
INSERT INTO posts (id, title, description, published_at, url, feed_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;


-- name: GetPostsForUser :many
SELECT posts.* from posts
JOIN feed_follows ON posts.feed_id = feed_follows.feed_id
WHERE feed_follows.user_id = $1 LIMIT $2;