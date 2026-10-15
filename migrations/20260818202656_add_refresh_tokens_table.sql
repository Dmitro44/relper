-- +goose Up
-- +goose StatementBegin
CREATE TABLE refresh_tokens (
    id uuid PRIMARY KEY default gen_random_uuid(),
    user_id uuid not null references users(id) on delete cascade,
    token_hash text not null unique,
    expires_at timestamptz not null,
    created_at timestamptz default now(),
    revoked_at timestamptz,
    replaced_by uuid
);

create index refresh_tokens_user_id_idx on refresh_tokens (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE refresh_tokens;
-- +goose StatementEnd
