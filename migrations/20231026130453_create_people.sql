-- +goose Up
-- +goose StatementBegin
CREATE TABLE people (
    id UUID NOT NULL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    middle_name TEXT,
    gender TEXT NOT NULL,
    nationality TEXT NOT NULL,
    age INT NOT NULL
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX people_first_name
    ON people (first_name);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE people;
-- +goose StatementEnd
