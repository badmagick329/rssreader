// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: feed_follows.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createFeedFollow = `-- name: CreateFeedFollow :one
INSERT INTO
  feed_follows (id, updated_at, created_at, user_id, feed_id)
VALUES
  ($1, $2, $3, $4, $5) RETURNING id, updated_at, created_at, user_id, feed_id
`

type CreateFeedFollowParams struct {
	ID        uuid.UUID
	UpdatedAt time.Time
	CreatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
}

func (q *Queries) CreateFeedFollow(ctx context.Context, arg CreateFeedFollowParams) (FeedFollow, error) {
	row := q.db.QueryRowContext(ctx, createFeedFollow,
		arg.ID,
		arg.UpdatedAt,
		arg.CreatedAt,
		arg.UserID,
		arg.FeedID,
	)
	var i FeedFollow
	err := row.Scan(
		&i.ID,
		&i.UpdatedAt,
		&i.CreatedAt,
		&i.UserID,
		&i.FeedID,
	)
	return i, err
}

const deleteFeedFollow = `-- name: DeleteFeedFollow :one
DELETE FROM feed_follows WHERE id=$1 RETURNING id, updated_at, created_at, user_id, feed_id
`

func (q *Queries) DeleteFeedFollow(ctx context.Context, id uuid.UUID) (FeedFollow, error) {
	row := q.db.QueryRowContext(ctx, deleteFeedFollow, id)
	var i FeedFollow
	err := row.Scan(
		&i.ID,
		&i.UpdatedAt,
		&i.CreatedAt,
		&i.UserID,
		&i.FeedID,
	)
	return i, err
}

const getFeedFollowByFeedID = `-- name: GetFeedFollowByFeedID :one
SELECT id, updated_at, created_at, user_id, feed_id FROM feed_follows WHERE feed_id=$1
`

func (q *Queries) GetFeedFollowByFeedID(ctx context.Context, feedID uuid.UUID) (FeedFollow, error) {
	row := q.db.QueryRowContext(ctx, getFeedFollowByFeedID, feedID)
	var i FeedFollow
	err := row.Scan(
		&i.ID,
		&i.UpdatedAt,
		&i.CreatedAt,
		&i.UserID,
		&i.FeedID,
	)
	return i, err
}

const getFeedFollows = `-- name: GetFeedFollows :many
SELECT id, updated_at, created_at, user_id, feed_id FROM feed_follows WHERE user_id=$1
`

func (q *Queries) GetFeedFollows(ctx context.Context, userID uuid.UUID) ([]FeedFollow, error) {
	rows, err := q.db.QueryContext(ctx, getFeedFollows, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FeedFollow
	for rows.Next() {
		var i FeedFollow
		if err := rows.Scan(
			&i.ID,
			&i.UpdatedAt,
			&i.CreatedAt,
			&i.UserID,
			&i.FeedID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
