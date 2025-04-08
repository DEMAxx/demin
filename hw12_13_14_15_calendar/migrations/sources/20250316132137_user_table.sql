-- +goose Up
-- +goose StatementBegin
CREATE TABLE "events"
(
    "id" UUID PRIMARY KEY,
    "name" TEXT NOT NULL,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "users";
-- +goose StatementEnd
