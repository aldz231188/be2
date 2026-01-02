-- name: CreateSession :exec
INSERT INTO sessions (jti_hash, user_id, expires_at)
VALUES ($1, $2, $3);

-- name: GetSessionByHash :one
SELECT jti_hash, user_id, expires_at, revoked_at, created_at
FROM sessions
WHERE jti_hash = $1
LIMIT 1;

-- name: RevokeSession :execrows
UPDATE sessions
SET revoked_at = now()
WHERE jti_hash = $1 AND revoked_at IS NULL;

-- name: RevokeSessionsByUser :execrows
UPDATE sessions
SET revoked_at = now()
WHERE user_id = $1 AND revoked_at IS NULL;
