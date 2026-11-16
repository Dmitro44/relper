-- +goose Up
-- +goose StatementBegin
-- rows without a location have no target left and can never satisfy the check below
DELETE FROM tips WHERE location_id IS NULL;

ALTER TABLE tips ADD COLUMN place_id uuid REFERENCES places(id) ON DELETE CASCADE;
ALTER TABLE tips ADD CONSTRAINT tips_one_target CHECK (num_nonnulls(location_id, place_id) = 1);

CREATE INDEX tips_place_idx ON tips (place_id);

ALTER TABLE comments ALTER COLUMN user_id SET NOT NULL;
ALTER TABLE tasks ALTER COLUMN user_id SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tasks ALTER COLUMN user_id DROP NOT NULL;
ALTER TABLE comments ALTER COLUMN user_id DROP NOT NULL;

DROP INDEX tips_place_idx;
ALTER TABLE tips DROP CONSTRAINT tips_one_target;
ALTER TABLE tips DROP COLUMN place_id;
-- +goose StatementEnd
