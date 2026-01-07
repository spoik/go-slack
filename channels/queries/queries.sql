-- name: ListChannels :many
SELECT * FROM channels;

-- name: CreateChannel :one
INSERT INTO channels (name)
VALUES ($1)
RETURNING *;
