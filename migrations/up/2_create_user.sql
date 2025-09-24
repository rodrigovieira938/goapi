CREATE TABLE IF NOT EXISTS "user" (
    id       SERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    email    TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);
