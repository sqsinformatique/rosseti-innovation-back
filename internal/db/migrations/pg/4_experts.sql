-- +goose Up
CREATE TABLE IF NOT EXISTS production.experts (
    id serial PRIMARY KEY,
    expert_id INTEGER NOT NULL,
    meta jsonb,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    deleted_at timestamp with time zone
);

-- +goose Down
DROP TABLE production.experts;
