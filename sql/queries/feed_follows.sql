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
    users.name AS user_name,
    feeds.name AS feed_name
FROM
    inserted_feed_follow
LEFT JOIN
    users ON inserted_feed_follow.user_id = users.id
LEFT JOIN
    feeds ON inserted_feed_follow.feed_id = feeds.id;

-- name: GetFeedFollowsForUser :many
SELECT feeds.name AS feed_name, users.name AS user_name
FROM feed_follows ff
JOIN feeds ON ff.feed_id = feeds.id
JOIN users ON ff.user_id = users.id
WHERE ff.user_id = $1;
