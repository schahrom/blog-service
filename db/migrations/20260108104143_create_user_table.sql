-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "user"
(
    user_id    serial,
    fist_name  varchar(64),
    last_name  varchar(64),
    login      varchar(64),
    email      varchar(64),
    password   varchar(64),
    created_at timestamp,

    CONSTRAINT PK_user_user_id PRIMARY KEY (user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "user"
    DROP CONSTRAINT PK_user_user_id;

DROP TABLE "user";
-- +goose StatementEnd
