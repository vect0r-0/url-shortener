-- +goose Up

CREATE TABLE  url (
    id UUID PRIMARY KEY,
    alias TEXT NOT NULL UNIQUE,
    url TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);



-- +goose Down

DROP TABLE url;

