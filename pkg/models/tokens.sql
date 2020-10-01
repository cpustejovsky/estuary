CREATE TABLE reset_tokens (
    id uuid DEFAULT,
    email VARCHAR NOT NULL,
    created timestamp without time zone default (now() at time zone 'utc'),
    PRIMARY KEY (id)
);