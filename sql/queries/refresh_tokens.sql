-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, expires_at, revoked_at, user_id)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    NUll, 
    $3
)
RETURNING *;

-- name: RevokeToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1;

-- name: GetRefreshToken :one
SELECT * 
FROM refresh_tokens 
WHERE token = $1;

-- name: GetUserFromRefreshToken :one
SELECT *
FROM users
WHERE id = (
    SELECT user_id
    FROM refresh_tokens
    WHERE token = $1
);