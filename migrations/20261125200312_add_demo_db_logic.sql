-- +goose Up
-- +goose StatementBegin
CREATE TRIGGER audit_places
AFTER INSERT OR UPDATE OR DELETE ON places
FOR EACH ROW EXECUTE FUNCTION audit_row_change();

CREATE TRIGGER audit_tasks
AFTER INSERT OR UPDATE OR DELETE ON tasks
FOR EACH ROW EXECUTE FUNCTION audit_row_change();

CREATE TRIGGER audit_likes
AFTER INSERT OR UPDATE OR DELETE ON likes
FOR EACH ROW EXECUTE FUNCTION audit_row_change();

CREATE OR REPLACE FUNCTION reject_self_like()
RETURNS trigger
LANGUAGE plpgsql
AS $$
DECLARE
    v_owner uuid;
BEGIN
    IF NEW.tip_id IS NOT NULL THEN
        SELECT user_id INTO v_owner FROM tips WHERE id = NEW.tip_id;
    ELSE
        SELECT user_id INTO v_owner FROM comments WHERE id = NEW.comment_id;
    END IF;

    IF v_owner = NEW.user_id THEN
        RAISE EXCEPTION 'users cannot like their own content'
            USING ERRCODE = '23514';
    END IF;

    RETURN NEW;
END;
$$;

CREATE TRIGGER prevent_self_like
BEFORE INSERT OR UPDATE ON likes
FOR EACH ROW EXECUTE FUNCTION reject_self_like();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS audit_places ON places;
DROP TRIGGER IF EXISTS audit_tasks ON tasks;
DROP TRIGGER IF EXISTS audit_likes ON likes;
DROP TRIGGER IF EXISTS prevent_self_like ON likes;
DROP FUNCTION IF EXISTS reject_self_like();
-- +goose StatementEnd
