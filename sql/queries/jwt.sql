-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens(token, user_id, expires_at)
VALUES (
  $1,
  $2,
  $3
)
RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT refresh_tokens.user_id, refresh_tokens.expires_at, refresh_tokens.revoked_at
FROM users
INNER JOIN refresh_tokens ON users.id = refresh_tokens.user_id
WHERE refresh_tokens.token = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET 
    revoked_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE token = $1;
