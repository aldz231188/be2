-- name: CreateUser :one
INSERT INTO client (id, client_name,client_surname,birthday,gender,address_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING registration_date;
-- name: AddAddress :exec
INSERT INTO address (id, country,city,street)
VALUES ($1, $2, $3, $4);


-- -- name: GetUserByNameAndSurename :one
-- SELECT *
-- FROM users
-- WHERE clien = $1;

-- -- name: ListUsers :many
-- SELECT id, login, password, created_at, deleted_at
-- FROM users
-- ORDER BY created_at DESC
-- LIMIT $1 OFFSET $2;

-- -- name: SoftDeleteUser :execrows
-- UPDATE users
-- SET deleted_at = now()
-- WHERE id = $1 AND deleted_at IS NULL;
