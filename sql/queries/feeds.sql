-- name: CreateFeed :one
insert into feeds (id, created_at, updated_at, name, url, user_id)
values($1, $2, $3, $4, $5, $6)
returning *;

-- name: GetFeeds :many
select * from feeds;

-- name: GetFeedsByUser :many
select feeds.id, feeds.created_at, feeds.updated_at, feeds.name, feeds.url, feeds.user_id, feeds.last_fetched_at
from feeds join feed_follows on feeds.id = feed_follows.feed_id
where feed_follows.user_id = $1;

-- name: GetNextFeedsToFetch :many
select *
from feeds
order by last_fetched_at asc nulls first
limit $1;

-- name: MarkFeedAsFetched :one
update feeds
set last_fetched_at = now(), updated_at = now()
where id = $1
returning *;
