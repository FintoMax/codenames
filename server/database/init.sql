DROP TABLE IF EXISTS users;
CREATE TABLE IF NOT EXISTS users
(
    id            SERIAL PRIMARY KEY,
    name          TEXT UNIQUE NOT NULL,
    password_hash TEXT        NOT NULL
);
