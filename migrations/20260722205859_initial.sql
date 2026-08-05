-- +goose Up
-- +goose StatementBegin
CREATE TABLE if not exists users (
    id uuid primary key default gen_random_uuid(),
    email text unique not null,
    nickname text unique not null,
    habits jsonb,
    contacts jsonb,
    created_at timestamptz default now()
);
CREATE TABLE IF NOT EXISTS locations (
    id uuid PRIMARY KEY default gen_random_uuid(),
    country text not null,
    city text not null,
    coordinates jsonb,
    created_at timestamptz default now()
);
CREATE TABLE IF NOT EXISTS places (
    id uuid PRIMARY KEY default gen_random_uuid(),
    location_id uuid references locations(id),
    name text not null,
    category text not null,
    address text,
    phone text,
    website_url text,
    coordinates jsonb,
    description text,
    rating numeric(5,1) default 0,
    upvotes int default 0,
    created_at timestamptz default now()
);
CREATE TABLE IF NOT EXISTS tasks (
    id uuid PRIMARY KEY default gen_random_uuid(),
    location_id uuid references locations(id),
    user_id uuid references users(id),
    title text not null,
    description text,
    place_id uuid references places(id),
    done boolean default false,
    created_at timestamptz default now()
);
CREATE TABLE IF NOT EXISTS tips (
    id uuid PRIMARY KEY default gen_random_uuid(),
    location_id uuid references locations(id),
    user_id uuid references users(id),
    category text not null,
    title text not null,
    content text not null,
    upvotes int default 0,
    created_at timestamptz default now()
);
CREATE TABLE IF NOT EXISTS tags (
    id uuid PRIMARY KEY default gen_random_uuid(),
    name text unique not null,
    slug text unique not null,
    created_at timestamptz default now()
);
CREATE TABLE IF NOT EXISTS places_tags (
    place_id uuid references places(id) on delete cascade,
    tag_id uuid references tags(id) on delete cascade,
    user_id uuid references users(id),
    created_at timestamptz default now(),
    PRIMARY KEY(place_id, tag_id, user_id)
);
CREATE TABLE IF NOT EXISTS comments (
    id uuid PRIMARY KEY default gen_random_uuid(),
    target_type text not null,            -- "place" / "tip" / "task"
    target_id uuid not null,
    user_id uuid references users(id),
    content text not null,
    created_at timestamptz default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS places_tags;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS tips;
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS places;
DROP TABLE IF EXISTS locations;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
