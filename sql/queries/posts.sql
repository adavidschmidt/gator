-- name: CreatePost :exec
INSERT INTO posts(id, created_at, updated_at, title, url, description, published_at, feed_id)
Values (
    $1,
    $2,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
);


-- name: GetPostsForUser :many
SELECT posts.*, feeds.name as feed_name
FROM posts
JOIN feed_follows ON feed_follows.feed_id = posts.feed_id
join feeds ON posts.feed_id = feeds.id
WHERE feeds.user_id = $1
LIMIT $2;