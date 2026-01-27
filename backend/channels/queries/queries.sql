-- name: ListChannels :many
SELECT *
FROM channels
ORDER BY name ASC;

-- name: CountChannels :one
SELECT COUNT(*) FROM channels;

-- name: CreateChannel :one
INSERT INTO channels (name)
VALUES ($1)
RETURNING *;

-- name: ChannelExists :one
SELECT EXISTS (
	SELECT 1 FROM channels WHERE id=$1
);

-- name: CreateMessage :one
INSERT INTO messages
(channel_id, message)
VALUES ($1, $2)
RETURNING *;

-- name: CountMessages :one
SELECT COUNT(*) FROM messages;

-- name: MessagesInChannel :many
SELECT *
FROM messages
WHERE channel_id = $1;
