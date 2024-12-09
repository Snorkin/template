create table if not exists public.client
(
    id    bigint generated always as identity
    primary key,
    name  text not null,
    email text not null,
    login text not null unique
);