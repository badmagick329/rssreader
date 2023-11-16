-- +goose Up
CREATE TABLE feed_follows (
  id UUID PRIMARY KEY,
  updated_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP NOT NULL,
  user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  feed_id UUID NOT NULL REFERENCES feeds (id) ON DELETE CASCADE
);
-- +goose Down
DROP TABLE feed_follows;