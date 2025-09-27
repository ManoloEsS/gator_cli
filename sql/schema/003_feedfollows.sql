-- +goose Up
CREATE TABLE feed_follows (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  user_id UUID NOT NULL,
  feed_id UUID NOT NULL,
  UNIQUE (user_id, feed_id),
  CONSTRAINT fk_user_id
  FOREIGN KEY (user_id)
  REFERENCES users(id)
  ON DELETE CASCADE,
  CONSTRAINT fk_feed_id
  FOREIGN KEY (feed_id)
  REFERENCES rssfeeds(id)
  ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS
idx_feed_follows_user_id
ON feed_follows(user_id);

CREATE INDEX IF NOT EXISTS
idx_feed_follows_feed_id
ON feed_follows(feed_id);

-- +goose Down
DROP INDEX IF EXISTS idx_feed_follows_user_id;
DROP INDEX IF EXISTS idx_feed_follows_feed_id;
DROP TABLE feed_follows;
