-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES ($1, NOW(), NOW(), $2,NOW() + INTERVAL '60 days', NULL)
RETURNING *;
--


-- name: GetUserFromRefreshToken :one
SELECT users.* FROM refresh_tokens
INNER JOIN users ON refresh_tokens.user_id = users.id
WHERE token = $1 AND expires_at > NOW() AND revoked_at IS NULL;
--

-- name: RevokeToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1;
--