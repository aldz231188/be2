CREATE TABLE IF NOT EXISTS "user" (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    login text UNIQUE NOT NULL,
    password_hash text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);
