-- name: CreatePost :one
insert into posts (id, created_at, updated_at, title, description, published_at, url, feed_id)
values ($1, $2, $3, $4, $5, $6, $7, $8)
returning *;

-- name: GetPostsForUser :many
select p.*
from posts p join feed_follows ff on p.feed_id = ff.feed_id
where ff.user_id = $1
order by p.published_at desc
limit $2;
