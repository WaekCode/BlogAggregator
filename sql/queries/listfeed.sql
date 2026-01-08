-- name: Listfeeds :many
SELECT feed.name,feed.url,COALESCE(users.name, '') AS user_name
FROM feed
LEFT JOIN users ON feed.user_id = users.id
;