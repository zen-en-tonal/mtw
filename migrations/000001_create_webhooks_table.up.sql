CREATE TABLE IF NOT EXISTS webhooks (
    id uuid,
    endpoint text,
    auth text,
    schema text,
    method text,
    content_type text,

    constraint webhooks_pk primary key (id)
);
