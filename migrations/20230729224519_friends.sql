-- +goose Up

-- +goose StatementBegin

CREATE TABLE
    friends (
        friend_one INT NOT NULL REFERENCES users (user_id),
        friend_two INT NOT NULL REFERENCES users (user_id)
    );

-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin

DROP TABLE IF EXISTS friends;

-- +goose StatementEnd