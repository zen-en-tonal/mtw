CREATE TABLE IF NOT EXISTS addresses_webhooks (
    address text,
    webhook_id uuid,

    constraint addresses_webhooks_pk primary key (address, webhook_id),
    foreign key (address) references addresses(address),
    foreign key (webhook_id) references webhooks(id)
);
