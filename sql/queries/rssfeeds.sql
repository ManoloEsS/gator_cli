-- name: CreateRSSFeed :one
INSERT INTO rssfeeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT users.name AS user_name, 
    JSON_AGG(JSON_BUILD_OBJECT('name', rssfeeds.name, 'url', rssfeeds.url)) AS feed_details 
FROM users 
INNER JOIN rssfeeds 
ON users.id = rssfeeds.user_id 
GROUP BY users.name 
ORDER BY users.name;

-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
  INSERT INTO feed_follows(id, created_at, updated_at, user_id, feed_id)
  VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
        )
  RETURNING *
)
SELECT
    inserted_feed_follow.*,
    rssfeeds.name AS feed_name,
    users.name AS user_name
FROM inserted_feed_follow
INNER JOIN users ON inserted_feed_follow.user_id = users.id
INNER JOIN rssfeeds ON inserted_feed_follow.feed_id = rssfeeds.id;

-- name: GetFeedByUrl :one
SELECT *
FROM rssfeeds
WHERE rssfeeds.Url = $1;
