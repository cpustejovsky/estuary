CREATE TABLE accounts (
    account_id uuid DEFAULT uuid_generate_v4 (),
    first_name VARCHAR NOT NULL,
    last_name VARCHAR NOT NULL,
    email VARCHAR NOT NULL UNIQUE,
    hashed_password CHAR(60) NOT NULL,
    created timestamp without time zone default (now() at time zone 'utc'),
    updated timestamp,
    deleted timestamp,
    email_updates BOOLEAN DEFAULT false,
    advanced_view BOOLEAN DEFAULT false,
    active BOOLEAN DEFAULT true,
    PRIMARY KEY (account_id)
);

INSERT INTO users (first_name, last_name, email, hashed_password, email_updates, advanced_view) VALUES (
    'Charles',
    'Pustejovsky',
    'charles.pustejovsky@gmail.com',
    '$2a$12$NuTjWXm3KKntReFwyBVHyuf/to.HEwTy.eS206TNfkGfr6HzGJSWG',
    true,
    true,
);
