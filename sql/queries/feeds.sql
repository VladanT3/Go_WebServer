-- name: CreateFeed :one
insert into feeds (id, created_at, updated_at, name, url, user_id)
values($1, $2, $3, $4, $5, $6)
returning *;

-- name: GetFeeds :many
select * from feeds;

-- name: GetFeedsByUser :many
select feeds.id, feeds.created_at, feeds.updated_at, feeds.name, feeds.url, feeds.user_id
from feeds join feed_follows on feeds.id = feed_follows.feed_id
where feed_follows.user_id = $1;
