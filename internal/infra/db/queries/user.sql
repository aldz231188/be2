-- name: GetUserByUsername :one
SELECT id, username, password_hash, created_at, token_version
FROM "user"
WHERE username = $1
LIMIT 1;

-- name: GetUserByID :one
SELECT id, username, password_hash, created_at, token_version
FROM "user"
WHERE id = $1
LIMIT 1;

-- name: CreateUser :one
INSERT INTO "user" (username, password_hash)
VALUES ($1, $2)
RETURNING id, username, password_hash, created_at, token_version;

-- name: IncrementTokenVersion :execrows
UPDATE "user"
SET token_version = token_version + 1
WHERE id = $1;
