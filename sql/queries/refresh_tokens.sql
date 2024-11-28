-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
  $1,
  NOW(),
  NOW(),
  $2,
  NOW() + INTERVAL '60' DAY,
  NULL
)
RETURNING *;

-- name: GetValidRefreshTokenForUser :one
SELECT token
FROM refresh_tokens
WHERE user_id = $1 AND
  revoked_at = NULL AND
  expires_at <= NOW()
ORDER BY created_at ASC
LIMIT 1;

-- name: GetUserFromRefreshToken :one
SELECT user_id
FROM refresh_tokens
WHERE token = $1 AND
  revoked_at IS NULL AND
  expires_at > NOW();

-- name: RevokeAccessToken :one
UPDATE refresh_tokens
SET updated_at = NOW(),
  revoked_at = NOW()
WHERE token = $1
RETURNING token, revoked_at;
