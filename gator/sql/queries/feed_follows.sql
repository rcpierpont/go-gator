-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
) SELECT 
    inserted_feed_follow.*,
    f.name AS feed_name,
    u.name AS user_name
FROM inserted_feed_follow
INNER JOIN users u ON u.id = inserted_feed_follow.user_id
INNER JOIN feeds f ON f.id = inserted_feed_follow.feed_id;

-- name: GetFeedFollowsForUser :many
WITH user_feed_follows AS (
    SELECT * FROM feed_follows WHERE feed_follows.user_id = $1
)
SELECT user_feed_follows.*,
    f.name AS feed_name,
    u.name AS user_name
FROM user_feed_follows
INNER JOIN users u ON u.id = user_feed_follows.user_id
INNER JOIN feeds f ON f.id = user_feed_follows.feed_id;

-- name: DeleteUserFeedFollow :exec
DELETE FROM feed_follows
WHERE feed_follows.user_id = $1 AND feed_follows.feed_id = $2;

