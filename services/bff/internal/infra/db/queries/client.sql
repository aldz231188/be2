-- name: CreateClient :exec
INSERT INTO public.client (id,client_name, client_surname, birthday, gender, address_id)
VALUES (
    sqlc.arg(id),
  sqlc.arg(client_name),
  sqlc.arg(client_surname),
  sqlc.arg(birthday),
  sqlc.arg(gender),
  sqlc.arg(address_id)
);;


-- name: UpdateClient :execrows
UPDATE public.client
SET client_name   = $2,
    client_surname= $3,
    birthday      = $4,
    gender        = $5,
    address_id    = $6
WHERE id = $1;

-- name: DeleteClient :execrows
DELETE FROM public.client
WHERE id = $1;
