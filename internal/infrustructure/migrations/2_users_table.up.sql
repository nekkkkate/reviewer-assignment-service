create table if not exists users (
    id serial primary key,
    name varchar(255) not null,
    email varchar(255) unique not null,
    team_name varchar(255) references teams(name) on delete cascade,
    is_active boolean default true not null
);