create table if not exists team_members (
    user_id int references users(id) on delete cascade,
    team_id int references teams(id) on delete cascade
);