create table if not exists prs (
    id serial primary key,
    title varchar(255) not null,
    author_id int references users(id) on delete cascade,
    team_id int references teams(id) on delete cascade,
    status varchar(255) default 'open' not null,
    created_at timestamp default current_timestamp not null,
    merged_at timestamp default null

);