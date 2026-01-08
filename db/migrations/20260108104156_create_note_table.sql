-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS note
(
    note_id    serial,
    title      text,
    content    text,
    create_at  timestamp,
    updated_at timestamp,
    user_id    int,

    CONSTRAINT PK_note_note_id PRIMARY KEY (note_id),
    CONSTRAINT FK_note_user_user_id FOREIGN KEY (user_id) REFERENCES "user" (user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE note
    DROP CONSTRAINT PK_note_note_id;
ALTER TABLE note
    DROP CONSTRAINT FK_note_user_user_id;
DROP TABLE note;
-- +goose StatementEnd
