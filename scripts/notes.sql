CREATE TABLE notes (
    note_id uuid DEFAULT uuid_generate_v4 (),
    content VARCHAR NOT NULL,
    category VARCHAR NOT NULL DEFAULT 'in-tray',
    tags text[],
    created timestamp without time zone default (now() at time zone 'utc'),
    due_date timestamp,
    remind_date timestamp,
    completedDate timestamp,
    account_id uuid references accounts(account_id),
    PRIMARY KEY (note_id)
);