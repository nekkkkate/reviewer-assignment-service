create index if not exists idx_users_team_name on users(team_name);
create index if not exists idx_users_is_active on users(is_active);
create index if not exists idx_users_email on users(email);

create index if not exists idx_teams_name on teams(name);

create index if not exists idx_team_members_user_id on team_members(user_id);
create index if not exists idx_team_members_team_id on team_members(team_id);
create unique index if not exists idx_team_members_unique on team_members(team_id, user_id);

create index if not exists idx_prs_author_id on prs(author_id);
create index if not exists idx_prs_team_id on prs(team_id);
create index if not exists idx_prs_status on prs(status);
create index if not exists idx_prs_created_at on prs(created_at);
create index if not exists idx_prs_merged_at on prs(merged_at);

create index if not exists idx_assigned_reviewers_pr_id on assigned_reviewers(pr_id);
create index if not exists idx_assigned_reviewers_user_id on assigned_reviewers(user_id);
create unique index if not exists idx_assigned_reviewers_unique on assigned_reviewers(pr_id, user_id);