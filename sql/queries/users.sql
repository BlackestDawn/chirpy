-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, password)
VALUES (
  gen_random_uuid(),
  NOW(),
  NOW(),
  $1,
  $2
)
RETURNING id, created_at, updated_at, email, is_chirpy_red;

-- name: ResetUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: UpdatePassword :one
UPDATE users
SET password = $1,
  email = $2,
  updated_at = NOW()
WHERE id = $3
RETURNING id, created_at, updated_at, email, is_chirpy_red;

-- name: UpgradeUserToRed :exec
UPDATE users
SET is_chirpy_red = TRUE
WHERE id = $1;
