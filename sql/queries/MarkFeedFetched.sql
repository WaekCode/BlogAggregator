-- name: MarkFeedFetched :exec
UPDATE feed
SET updated_at = NOW(),
    last_fetched_at = NOW()
WHERE id = $1;
