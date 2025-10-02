
-- name: CreateAddress :exec
INSERT INTO address (id, country,city,street)
VALUES ($1, $2, $3, $4);

-- name: UpdateAddress :execrows
UPDATE public.address
SET country = $2,
    city    = $3,
    street  = $4
WHERE id = $1;

-- name: DeleteAddress :execrows
DELETE FROM public.address
WHERE id = $1;

