create table if not exists assigned_reviewers (
    pr_id int references prs(id) on delete cascade,
    user_id int references users(id) on delete cascade
);