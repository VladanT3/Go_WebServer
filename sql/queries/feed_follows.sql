-- name: CreateFeedFollow :one
insert into feed_follows(id, created_at, updated_at, user_id, feed_id)
values($1, $2, $3, $4, $5)
returning *;

-- name: DeleteFeedFollow :exec
delete from feed_follows where user_id = $1 and feed_id = $2;
