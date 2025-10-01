-- +goose Up
ALTER TABLE rssfeeds
ADD COLUMN last_fetched_at TIMESTAMP;

-- +goose Down
ALTER TABLE rssfeeds
DROP COLUMN last_fetched_at;
