-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetFeedsByUser :many

SELECT * from feeds where user_id = $1;

-- name: GetFeeds :many

SELECT * from feeds;

-- name: CreateFeedFollow :one
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetFeedFollows :many

SELECT * from feed_follows;

-- name: GetFeedFollowsByUser :many

SELECT * from feed_follows where user_id=$1;

-- name: GetFeedFollowsByFeed :many

SELECT * from feed_follows where feed_id=$1;

-- name: DeleteFeedFollow :exec

DELETE FROM feed_follows WHERE id = $1 AND user_id = $2;