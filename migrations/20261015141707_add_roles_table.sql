-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS roles (
    id uuid PRIMARY KEY default gen_random_uuid(),
    code text unique not null,
    name text not null,
    created_at timestamptz default now()
);

CREATE TABLE IF NOT EXISTS user_roles (
    user_id uuid references users(id) on delete cascade,
    role_id uuid references roles(id) on delete cascade,
    PRIMARY KEY(user_id, role_id)
);

INSERT INTO roles (code, name)
VALUES ('admin', 'Administrator'), ('user', 'Regular user')
ON CONFLICT (code) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS roles;
-- +goose StatementEnd
