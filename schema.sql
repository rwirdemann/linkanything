CREATE TABLE links
(
    id      serial PRIMARY KEY,
    title   VARCHAR(255) NOT NULL,
    uri     VARCHAR(255) NOT NULL,
    tags    VARCHAR(500)          default '',
    draft   BOOLEAN               DEFAULT FALSE,
    created TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE users
(
    id       serial PRIMARY KEY,
    name     VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    created  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE tags
(
    id   serial PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE tags_links
(
    tag_id  integer,
    link_id integer
);
