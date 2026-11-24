-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_log (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    source text NOT NULL CHECK (source IN ('app', 'trigger')),
    action text NOT NULL,
    entity_type text NOT NULL,
    entity_id uuid,
    actor_id uuid REFERENCES users(id) ON DELETE SET NULL,
    details jsonb NOT NULL DEFAULT '{}'::jsonb,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX audit_log_created_at_idx ON audit_log (created_at DESC);
CREATE INDEX audit_log_actor_id_idx ON audit_log (actor_id);
CREATE INDEX audit_log_entity_idx ON audit_log (entity_type, entity_id);
CREATE INDEX audit_log_source_idx ON audit_log (source);

CREATE OR REPLACE FUNCTION audit_row_change()
RETURNS trigger
LANGUAGE plpgsql
AS $$
DECLARE
    v_actor uuid;
    v_entity uuid;
    v_details jsonb;
BEGIN
    v_entity := COALESCE((to_jsonb(NEW)->>'id')::uuid, (to_jsonb(OLD)->>'id')::uuid);
    v_actor := COALESCE((to_jsonb(NEW)->>'user_id')::uuid, (to_jsonb(OLD)->>'user_id')::uuid);
    v_details := jsonb_build_object(
        'old', CASE WHEN TG_OP IN ('UPDATE', 'DELETE') THEN to_jsonb(OLD) ELSE NULL END,
        'new', CASE WHEN TG_OP IN ('INSERT', 'UPDATE') THEN to_jsonb(NEW) ELSE NULL END
    );

    INSERT INTO audit_log (source, action, entity_type, entity_id, actor_id, details)
    VALUES ('trigger', lower(TG_TABLE_NAME) || '.' || lower(TG_OP), TG_TABLE_NAME, v_entity, v_actor, v_details);

    IF TG_OP = 'DELETE' THEN
        RETURN OLD;
    END IF;
    RETURN NEW;
END;
$$;

CREATE TRIGGER audit_place_ratings
AFTER INSERT OR UPDATE OR DELETE ON place_ratings
FOR EACH ROW EXECUTE FUNCTION audit_row_change();

CREATE TRIGGER audit_comments
AFTER INSERT OR UPDATE OR DELETE ON comments
FOR EACH ROW EXECUTE FUNCTION audit_row_change();

CREATE TRIGGER audit_tips
AFTER INSERT OR UPDATE OR DELETE ON tips
FOR EACH ROW EXECUTE FUNCTION audit_row_change();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS audit_place_ratings ON place_ratings;
DROP TRIGGER IF EXISTS audit_comments ON comments;
DROP TRIGGER IF EXISTS audit_tips ON tips;
DROP FUNCTION IF EXISTS audit_row_change();
DROP TABLE IF EXISTS audit_log;
-- +goose StatementEnd
