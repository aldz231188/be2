ALTER TABLE "user"
    ADD COLUMN IF NOT EXISTS token_version integer NOT NULL DEFAULT 1;

CREATE TABLE IF NOT EXISTS sessions (
    jti_hash text PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    expires_at timestamptz NOT NULL,
    revoked_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS sessions_user_id_idx ON sessions(user_id);
