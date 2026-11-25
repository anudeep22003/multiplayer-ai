-- Finds an existing user or creates a new one in a single database call,
-- and returns whether the user was just created (is_new).
--
-- 1. ON CONFLICT (username) DO UPDATE SET username = EXCLUDED.username
--    Performs a no-op update on conflict so RETURNING hands back the existing row.
-- 2. (xmax = 0)::boolean AS is_new
--    Freshly inserted rows have xmax = 0 (true); updated rows have non-zero xmax (false).
-- name: UpsertUser :one
INSERT INTO users (username) VALUES ($1)
ON CONFLICT (username) DO UPDATE SET username = EXCLUDED.username
RETURNING id, username, (xmax = 0)::boolean AS is_new;
