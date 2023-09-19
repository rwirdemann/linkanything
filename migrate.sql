CREATE TABLE links
(
    id       serial PRIMARY KEY,
    title    VARCHAR(255) NOT NULL,
    uri      VARCHAR(255) NOT NULL,
    tags     VARCHAR(500),
    created  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);