-- Write your migrate up statements here

create table sessions (
    token text primary key,
    data bytea not null,
    expiry timestamptz not null
);

create index sessions_expiry_idx on sessions (expiry);

---- create above / drop below ----
drop index if exists sessions_expiry_idx;
drop table if exists sessions;
-- Write your migrate down statements here. 
-- If this migration is irreversible Then delete the separator line above.
