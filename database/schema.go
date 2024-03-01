package datastore

const Schema = `
CREATE TABLE IF NOT EXISTS webhooks (
    id uuid,
    endpoint text,
    auth text,
    schema text,
    method text,
    content_type text,

    constraint webhooks_pk primary key (id)
);

CREATE TABLE IF NOT EXISTS addresses_webhooks (
    address text,
    webhook_id uuid,

    constraint addresses_webhooks_pk primary key (address, webhook_id)
);
`
