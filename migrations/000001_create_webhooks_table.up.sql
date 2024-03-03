CREATE TABLE IF NOT EXISTS webhooks (
    id uuid NOT NULL,
    endpoint text NOT NULL,
    auth text NOT NULL,
    schema text NOT NULL,
    method text NOT NULL,
    content_type text NOT NULL,

    constraint webhooks_pk primary key (id)
);
