-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password) 
VALUES (
  $1,
  $2,
  $3,
  $4,
  $5
)
RETURNING *;


-- name: DeleteUsers :exec

DELETE FROM users;

-- name: GetPass :one

SELECT * FROM users 
where email = $1;

-- name: UpdatePass :one
UPDATE users
  SET hashed_password = $1, email = $2, updated_at = $3
  WHERE id = $4
RETURNING *;


-- name: UpgradeChirpy :one
UPDATE users
SET chirpy_red = $1 
where id = $2
RETURNING *;
