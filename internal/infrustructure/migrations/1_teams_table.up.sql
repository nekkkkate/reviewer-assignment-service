create table if not exists teams (
    id serial primary key,
    name varchar(255) unique not null
)