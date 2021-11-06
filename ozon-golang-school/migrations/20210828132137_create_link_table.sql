-- +goose Up
CREATE TABLE IF NOT EXISTS links
(
    id          bigserial primary key,
    user_id     int8        not null,
    url         text        not null,
    description text        not null,
    tags        text        not null,
    created_at  timestamptz not null default now(),
    updated_at  timestamptz not null default current_timestamp
);

CREATE INDEX IF NOT EXISTS idx_user ON links USING hash (user_id);

-- +goose Down
DROP INDEX IF EXISTS idx_user;
DROP TABLE IF EXISTS links;