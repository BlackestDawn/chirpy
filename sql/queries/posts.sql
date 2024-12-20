-- name: CreatePost :one
INSERT INTO posts(id, created_at, updated_at, body, user_id)
VALUES (
  gen_random_uuid(),
  NOW(),
  NOW(),
  $1,
  $2
)
returning *;

-- name: ListPosts :many
SELECT *
FROM posts
ORDER BY created_at ASC;

-- name: GetPostByID :one
SELECT *
FROM posts
WHERE id = $1;

-- name: DeletePostByID :exec
DELETE FROM posts
WHERE id = $1;

-- name: ListPostsFromUser :many
SELECT *
FROM posts
WHERE user_id = $1
ORDER BY created_at ASC;
