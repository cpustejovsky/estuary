CREATE TABLE reset_tokens (
    token_id uuid NOT NULL,
    email VARCHAR NOT NULL,
    created timestamp without time zone default (now() at time zone 'utc'),
    PRIMARY KEY (token_id)
);