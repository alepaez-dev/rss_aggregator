-- +goose Up

ALTER TABLE users
ALTER COLUMN created_at SET DEFAULT now(),
ALTER COLUMN updated_at SET DEFAULT now();

ALTER TABLE feeds
ALTER COLUMN created_at SET DEFAULT now(),
ALTER COLUMN updated_at SET DEFAULT now();

ALTER TABLE feed_follows
ALTER COLUMN created_at SET DEFAULT now(),
ALTER COLUMN updated_at SET DEFAULT now();
-- +goose Down

-- Remove defaults (rollback)
ALTER TABLE users
ALTER COLUMN created_at DROP DEFAULT,
ALTER COLUMN updated_at DROP DEFAULT;

ALTER TABLE feeds
ALTER COLUMN created_at DROP DEFAULT,
ALTER COLUMN updated_at DROP DEFAULT;

ALTER TABLE feeds_follows
ALTER COLUMN created_at DROP DEFAULT,
ALTER COLUMN updated_at DROP DEFAULT;