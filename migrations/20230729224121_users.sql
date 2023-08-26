-- +goose Up

-- +goose StatementBegin

CREATE TABLE
    users (
        user_id INT NOT NULL UNIQUE,
        username VARCHAR(45) NOT NULL,
        age VARCHAR(3) NOT NULL
    );

-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin

DROP TABLE IF EXISTS users;

-- +goose StatementEnd