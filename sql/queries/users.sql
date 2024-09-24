-- name: CreateUser :one
INSERT INTO public.chirps (id, created_at, updated_at, email)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1
)
RETURNING *;