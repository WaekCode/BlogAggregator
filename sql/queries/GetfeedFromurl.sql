-- name: GetFeedFromUrl :one
SELECT * FROM feed
WHERE url = $1;