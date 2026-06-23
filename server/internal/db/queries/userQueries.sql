-- name: CreateUsersTable :exec
CREATE TABLE IF NOT EXISTS users (
    user_id BIGSERIAL PRIMARY KEY,
    email VARCHAR NOT NULL UNIQUE,
    password VARCHAR NOT NULL,
    refresh_token VARCHAR,
    creation_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- name: CreateUser :one
INSERT INTO users (email, password)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateRefreshToken :exec
UPDATE users
  SET refresh_token = $1
  WHERE user_id = $2;

-- name: GetRefreshTokenByUserID :one
SELECT refresh_token
FROM users
WHERE user_id = $1;  

-- name: SetRefreshTokenToNULL :exec 

UPDATE users
  SET refresh_token = NULL
  WHERE user_id = $1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: DeleteUserByID :exec
DELETE FROM users
WHERE user_id = $1;

-- name: DeleteUserByEmail :exec
DELETE FROM users
WHERE email = $1;

-- name: ListUsers :many
SELECT *
FROM users
ORDER BY user_id;
