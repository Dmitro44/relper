-- +goose Up
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN habits, DROP COLUMN contacts;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN habits jsonb, ADD COLUMN contacts jsonb;
-- +goose StatementEnd
