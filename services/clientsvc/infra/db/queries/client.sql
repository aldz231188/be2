-- name: CreateClient :exec
INSERT INTO client (id,user_id,client_name, client_surname, birthday, gender, address_id)
VALUES (
  sqlc.arg(id),
  sqlc.arg(user_id),
  sqlc.arg(client_name),
  sqlc.arg(client_surname),
  sqlc.arg(birthday),
  sqlc.arg(gender),
  sqlc.arg(address_id)
);;


-- name: UpdateClient :execrows
UPDATE client
SET user_id=$2,
    client_name   = $3,
    client_surname= $4,
    birthday      = $5,
    gender        = $6,
    address_id    = $7
WHERE id = $1;

-- name: DeleteClient :execrows
DELETE FROM client
WHERE id = $1;
