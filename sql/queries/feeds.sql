-- name: CreateFeed :one
INSERT INTO
  feeds (id, created_at, updated_at, name, url, user_id)
VALUES
  (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
  ) RETURNING *;

-- name: GetFeeds :many
SELECT
  *
FROM
  feeds
WHERE
  user_id = $1;

-- name: GetFeedsAll :many
SELECT
  *
FROM
  feeds;

-- name: GetFeed :one
SELECT
  *
FROM
  feeds
WHERE
  url = $1;

-- name: RemoveAllFeeds :exec
DELETE FROM
  feeds;
