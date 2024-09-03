-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, description, published_at, url, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPostsByFeedId :many

SELECT * from posts where feed_id = $1;

-- name: GetAllPosts :many

SELECT * FROM posts;

-- name: GetPostsForUser :many

SELECT p.* from posts p
JOIN feed_follows ff on p.feed_id=ff.feed_id
WHERE ff.user_id = $1
ORDER BY p.published_at DESC
LIMIT $2;