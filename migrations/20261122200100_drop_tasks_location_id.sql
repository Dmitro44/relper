-- +goose Up
-- +goose StatementBegin
ALTER TABLE tasks DROP COLUMN location_id;
ALTER TABLE tasks ALTER COLUMN place_id SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tasks ALTER COLUMN place_id DROP NOT NULL;
ALTER TABLE tasks ADD COLUMN location_id uuid REFERENCES locations(id);
-- +goose StatementEnd
