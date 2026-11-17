-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS place_ratings (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    place_id   uuid NOT NULL REFERENCES places(id) ON DELETE CASCADE,
    stars      smallint NOT NULL CHECK (stars BETWEEN 1 AND 5),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (user_id, place_id)
);

CREATE INDEX IF NOT EXISTS place_ratings_place_idx ON place_ratings (place_id);

CREATE TABLE IF NOT EXISTS likes (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tip_id     uuid REFERENCES tips(id) ON DELETE CASCADE,
    comment_id uuid REFERENCES comments(id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT now(),
    CHECK (num_nonnulls(tip_id, comment_id) = 1)
);

CREATE UNIQUE INDEX IF NOT EXISTS likes_user_tip_uq ON likes (user_id, tip_id) WHERE tip_id IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS likes_user_comment_uq ON likes (user_id, comment_id) WHERE comment_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS likes_tip_idx ON likes (tip_id) WHERE tip_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS likes_comment_idx ON likes (comment_id) WHERE comment_id IS NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS likes_comment_idx;
DROP INDEX IF EXISTS likes_tip_idx;
DROP INDEX IF EXISTS likes_user_comment_uq;
DROP INDEX IF EXISTS likes_user_tip_uq;
DROP TABLE IF EXISTS likes;

DROP INDEX IF EXISTS place_ratings_place_idx;
DROP TABLE IF EXISTS place_ratings;
-- +goose StatementEnd
