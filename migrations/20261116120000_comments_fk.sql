-- +goose Up
-- +goose StatementBegin
ALTER TABLE comments ADD COLUMN place_id uuid REFERENCES places(id) ON DELETE CASCADE;
ALTER TABLE comments ADD COLUMN tip_id uuid REFERENCES tips(id) ON DELETE CASCADE;

DELETE FROM comments c
WHERE c.target_type NOT IN ('place', 'tip')
   OR (c.target_type = 'place' AND NOT EXISTS (SELECT 1 FROM places p WHERE p.id = c.target_id))
   OR (c.target_type = 'tip' AND NOT EXISTS (SELECT 1 FROM tips t WHERE t.id = c.target_id));

UPDATE comments SET place_id = target_id WHERE target_type = 'place';
UPDATE comments SET tip_id = target_id WHERE target_type = 'tip';

ALTER TABLE comments ADD CONSTRAINT comments_one_target CHECK (num_nonnulls(place_id, tip_id) = 1);

ALTER TABLE comments DROP COLUMN target_type;
ALTER TABLE comments DROP COLUMN target_id;

CREATE INDEX comments_place_idx ON comments (place_id, created_at) WHERE place_id IS NOT NULL;
CREATE INDEX comments_tip_idx ON comments (tip_id, created_at) WHERE tip_id IS NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE comments ADD COLUMN target_type text;
ALTER TABLE comments ADD COLUMN target_id uuid;

UPDATE comments SET target_type = 'place', target_id = place_id WHERE place_id IS NOT NULL;
UPDATE comments SET target_type = 'tip', target_id = tip_id WHERE tip_id IS NOT NULL;

ALTER TABLE comments ALTER COLUMN target_type SET NOT NULL;
ALTER TABLE comments ALTER COLUMN target_id SET NOT NULL;

ALTER TABLE comments DROP CONSTRAINT comments_one_target;
DROP INDEX comments_place_idx;
DROP INDEX comments_tip_idx;

ALTER TABLE comments DROP COLUMN place_id;
ALTER TABLE comments DROP COLUMN tip_id;
-- +goose StatementEnd
