-- +goose Up


CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    full_name varchar(128) not null,
    email varchar(128) unique not null,
    password text not null,
    created_at timestamp DEFAULT NOW(),
    updated_at timestamp DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS users;
