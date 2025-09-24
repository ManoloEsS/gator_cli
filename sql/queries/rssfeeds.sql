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

