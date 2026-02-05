-- name: GetUserByLogin :one
SELECT id, login, password_hash, created_at, token_version
FROM "user"
WHERE login = $1
LIMIT 1;

-- name: GetUserByID :one
SELECT id, login, password_hash, created_at, token_version
FROM "user"
WHERE id = $1
LIMIT 1;

-- name: CreateUser :one
INSERT INTO "user" (login, password_hash)
VALUES ($1, $2)
RETURNING id, login, password_hash, created_at, token_version;

-- name: IncrementTokenVersion :execrows
UPDATE "user"
SET token_version = token_version + 1
WHERE id = $1;
