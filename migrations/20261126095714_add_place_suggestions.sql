-- +goose Up
-- +goose StatementBegin
CREATE TABLE place_suggestions (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    location_id uuid NOT NULL REFERENCES locations(id) ON DELETE CASCADE,
    name text NOT NULL,
    category text NOT NULL,
    address text NOT NULL DEFAULT '',
    phone text NOT NULL DEFAULT '',
    website_url text NOT NULL DEFAULT '',
    description text NOT NULL DEFAULT '',
    latitude double precision NOT NULL CHECK (latitude BETWEEN -90 AND 90),
    longitude double precision NOT NULL CHECK (longitude BETWEEN -180 AND 180),
    status text NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    rejection_reason text,
    reviewed_by uuid REFERENCES users(id) ON DELETE SET NULL,
    reviewed_at timestamptz,
    place_id uuid REFERENCES places(id) ON DELETE SET NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX place_suggestions_status_created_at_idx ON place_suggestions (status, created_at DESC);
CREATE INDEX place_suggestions_user_id_idx ON place_suggestions (user_id);
CREATE INDEX place_suggestions_location_id_idx ON place_suggestions (location_id);

CREATE TRIGGER audit_place_suggestions
AFTER INSERT OR UPDATE OR DELETE ON place_suggestions
FOR EACH ROW EXECUTE FUNCTION audit_row_change();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS audit_place_suggestions ON place_suggestions;
DROP TABLE IF EXISTS place_suggestions;
-- +goose StatementEnd
