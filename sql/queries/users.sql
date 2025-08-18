-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING id, created_at, updated_at, email;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserIDByRefreshToken :one
SELECT users.id
FROM users
INNER JOIN refresh_tokens ON users.id = refresh_tokens.user_id
WHERE refresh_tokens.token = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: UpdateUser :exec
UPDATE users
SET email = $2, hashed_password = $3
WHERE id = $1;

-- name: UpgradeToRed :exec
UPDATE users
SET is_chirpy_red = TRUE
WHERE id = $1;

-- name: RevokeRed :exec
UPDATE users
SET is_chirpy_red = FALSE
WHERE id = $1;