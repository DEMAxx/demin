-- +goose Up
-- +goose StatementBegin
CREATE TABLE "events"
(
    "id" UUID PRIMARY KEY,
    "title" TEXT NOT NULL,
    "date" TIMESTAMP NOT NULL,
    "duration" INTERVAL NOT NULL,
    "description" TEXT,
    "user_id" UUID NOT NULL,
    "notify" TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "events";
-- +goose StatementEnd
