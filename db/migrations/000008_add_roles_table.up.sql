CREATE TABLE IF NOT EXISTS roles (
    id bigserial PRIMARY KEY,
    name varchar(255) NOT NULL UNIQUE,
    level int NOT NULL DEFAULT 0,
    description TEXT
);