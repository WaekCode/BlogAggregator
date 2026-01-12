-- +goose Up
CREATE TABLE posts (
    id UUID PRIMARY KEY,                  -- unique identifier
    created_at TIMESTAMP NOT NULL,        -- time record created
    updated_at TIMESTAMP NOT NULL,        -- time record last updated
    title TEXT NOT NULL,                  -- title of the post
    url TEXT NOT NULL UNIQUE,             -- URL of the post, must be unique
    description TEXT NOT NULL,            -- description of the post
    published_at TIMESTAMP NOT NULL,      -- time post was published
    feed_id UUID NOT NULL REFERENCES feed(id)  -- foreign key to feed
);




-- +goose Down
DROP TABLE posts;