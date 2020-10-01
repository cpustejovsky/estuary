CREATE TABLE password_reset_tokens (
    id uuid DEFAULT uuid_generate_v4 (),
    email VARCHAR NOT NULL,
    created timestamp without time zone default (now() at time zone 'utc'),
    PRIMARY KEY (id)
);