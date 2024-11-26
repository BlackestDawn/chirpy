-- name: CreatePost :one
INSERT INTO posts(id, created_at, updated_at, body, user_id)
VALUES ($1, $2, $3, $4, $5)
returning *;

-- name: ListPosts :many
SELECT *
FROM posts
ORDER BY created_at ASC;

-- name: GetPostByID :one
SELECT *
FROM posts
WHERE id = $1;
