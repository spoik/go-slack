-- name: ListChannels :many
SELECT *
FROM channels
ORDER BY name ASC;

-- name: CreateChannel :one
INSERT INTO channels (name)
VALUES ($1)
RETURNING *;

-- name: ChannelExists :one
SELECT EXISTS (
	SELECT 1 FROM channels WHERE id=$1
);
