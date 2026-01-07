-- name: ListChannels :many
SELECT *
FROM channels
ORDER BY name ASC;

-- name: CreateChannel :one
INSERT INTO channels (name)
VALUES ($1)
RETURNING *;
