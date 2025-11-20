-- +goose Up

ALTER TABLE users ADD COLUMN api_key VARCHAR(64) UNIQUE NOT NULL DEFAULT(
    encode(sha256(random()::text::bytea), 'hex')
);

-- +goose Down
ALTER TABLE users DROP COLUMN api_key;

--goose postgres postgres://postgres:HC18@localhost:5432/RSS-Aggregator up