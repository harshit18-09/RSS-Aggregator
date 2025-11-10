-- name: CreateUser :one
INSERT INTO users (id, name, created_at, updated_at)
VALUES ($1, $2, $3, $4)                                --$number is a placeholder for parameters like $1 for id etc
RETURNING *;
