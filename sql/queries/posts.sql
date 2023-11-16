-- name: CreatePost :one
INSERT INTO
  posts (id, created_at, updated_at, feed_id, title, url, description, published_at)
VALUES
  (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
  ) RETURNING *;
-- name: GetPostsByUserID :many
SELECT
  *
FROM
  posts
WHERE
  feed_id IN (
    SELECT
      id
    FROM
      feeds
    WHERE
      user_id = $1
  );
