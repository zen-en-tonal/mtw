package datastore

const Schema = `
CREATE TABLE IF NOT EXISTS webhooks (
    id uuid,
    endpoint text,
    auth text,
    schema text,
    method text,
    content_type text
);

CREATE TABLE IF NOT EXISTS addresses_webhooks (
    address text,
    webhook_id uuid
);
`
