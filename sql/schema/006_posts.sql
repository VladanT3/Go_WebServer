-- +goose Up
CREATE TABLE posts(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title text not null,
    description text,
    published_at timestamp not null,
    url text not null unique,
    feed_id uuid not null references feeds(id) on delete cascade
);
-- +goose Down
DROP TABLE posts;
