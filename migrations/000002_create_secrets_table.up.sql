CREATE TABLE IF NOT EXISTS secrets (
    secret_id uuid DEFAULT gen_random_uuid () primary key,
    user_id   uuid REFERENCES users (user_id) on delete cascade,
    name           varchar(256) not null,
    kind           varchar(128) not null,
    metadata       bytea,
    data           bytea not null
)
