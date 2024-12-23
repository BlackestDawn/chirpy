// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: refresh_tokens.sql

package database

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createRefreshToken = `-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
  $1,
  NOW(),
  NOW(),
  $2,
  NOW() + INTERVAL '60' DAY,
  NULL
)
RETURNING token, created_at, updated_at, user_id, expires_at, revoked_at
`

type CreateRefreshTokenParams struct {
	Token  string    `json:"token"`
	UserID uuid.UUID `json:"user_id"`
}

func (q *Queries) CreateRefreshToken(ctx context.Context, arg CreateRefreshTokenParams) (RefreshToken, error) {
	row := q.db.QueryRowContext(ctx, createRefreshToken, arg.Token, arg.UserID)
	var i RefreshToken
	err := row.Scan(
		&i.Token,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.ExpiresAt,
		&i.RevokedAt,
	)
	return i, err
}

const getUserFromRefreshToken = `-- name: GetUserFromRefreshToken :one
SELECT user_id
FROM refresh_tokens
WHERE token = $1 AND
  revoked_at IS NULL AND
  expires_at > NOW()
`

func (q *Queries) GetUserFromRefreshToken(ctx context.Context, token string) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, getUserFromRefreshToken, token)
	var user_id uuid.UUID
	err := row.Scan(&user_id)
	return user_id, err
}

const getValidRefreshTokenForUser = `-- name: GetValidRefreshTokenForUser :one
SELECT token
FROM refresh_tokens
WHERE user_id = $1 AND
  revoked_at = NULL AND
  expires_at <= NOW()
ORDER BY created_at ASC
LIMIT 1
`

func (q *Queries) GetValidRefreshTokenForUser(ctx context.Context, userID uuid.UUID) (string, error) {
	row := q.db.QueryRowContext(ctx, getValidRefreshTokenForUser, userID)
	var token string
	err := row.Scan(&token)
	return token, err
}

const revokeRefreshToken = `-- name: RevokeRefreshToken :one
UPDATE refresh_tokens
SET updated_at = NOW(),
  revoked_at = NOW()
WHERE token = $1
RETURNING token, revoked_at
`

type RevokeRefreshTokenRow struct {
	Token     string       `json:"token"`
	RevokedAt sql.NullTime `json:"revoked_at"`
}

func (q *Queries) RevokeRefreshToken(ctx context.Context, token string) (RevokeRefreshTokenRow, error) {
	row := q.db.QueryRowContext(ctx, revokeRefreshToken, token)
	var i RevokeRefreshTokenRow
	err := row.Scan(&i.Token, &i.RevokedAt)
	return i, err
}
