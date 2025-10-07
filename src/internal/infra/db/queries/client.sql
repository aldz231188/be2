-- name: CreateClient :exec
INSERT INTO public.client (id,client_name, client_surname, birthday, gender, address_id)
VALUES (
  $1,
  $2,                     -- client_name
  $3,                     -- client_surname
  $4,                     -- birthday (date, CHECK birthday <= CURRENT_DATE)
  COALESCE($5, 'unknown')::public.gender_t,  -- gender (nullable param -> default 'unknown')
  $6                      -- address_id (nullable)
);

-- name: UpdateClient :execrows
UPDATE public.client
SET client_name   = $2,
    client_surname= $3,
    birthday      = $4,
    gender        = COALESCE($5, gender), -- если не хочешь менять, передай NULL
    address_id    = $6
WHERE id = $1;

-- name: DeleteClient :execrows
DELETE FROM public.client
WHERE id = $1;
