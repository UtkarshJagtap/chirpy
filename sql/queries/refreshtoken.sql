-- name: CreateRefreshToken :one
INSERT INTO refresh_token (
  token, created_at, updated_at, user_id, expires_at, revoked_at
) VALUES ( $1, $2, $3, $4, $5, $6)
returning token;

-- name: GetUserFromRefreshToken :one
SELECT user_id from refresh_token 
where token = $1 and revoked_at is null and expires_at > Now();

-- name: RevokeRefreshToken :exec
UPDATE refresh_token
  SET updated_at = $2, revoked_at = $2
  WHERE token = $1;

