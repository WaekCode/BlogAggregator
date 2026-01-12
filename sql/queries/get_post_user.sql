-- name: GetPostForUser :many
SELECT posts.* FROM posts

LEFT JOIN feed 
ON posts.feed_id = feed.id

WHERE feed.user_id = $1

ORDER BY posts.created_at DESC 
LIMIT $2
;