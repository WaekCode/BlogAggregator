-- +goose Up
ALTER TABLE feed
ADD COLUMN last_fetched_at TIMESTAMP NULL;




-- +goose Down
ALTER TABLE feed
DROP COLUMN last_fetched_at;

