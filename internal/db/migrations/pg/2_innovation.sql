-- +goose Up
CREATE TABLE IF NOT EXISTS production.innovation (
    id serial PRIMARY KEY,
    author_id INTEGER NOT NULL,
    title character varying(255) DEFAULT '',
    tags character varying(255) DEFAULT '',
	descriptions text DEFAULT '',
    state character varying(255) DEFAULT '',
    meta jsonb,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    deleted_at timestamp with time zone
);

-- +goose Down
DROP TABLE production.innovation;
