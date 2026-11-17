-- +goose Up
-- +goose StatementBegin
ALTER TABLE places DROP COLUMN IF EXISTS rating;
ALTER TABLE tips DROP COLUMN IF EXISTS upvotes;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE places ADD COLUMN rating numeric;
ALTER TABLE tips ADD COLUMN upvotes int default 0;
-- +goose StatementEnd
