-- name: GetFeedFollowsFromUser :many
SELECT
    feed_follows.*,
    users.name AS user_name,
    feed.name AS feed_name
FROM feed_follows
JOIN users ON users.id = feed_follows.user_id
JOIN feed ON feed.id = feed_follows.feed_id
WHERE feed_follows.user_id = $1;
