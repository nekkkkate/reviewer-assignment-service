drop index if exists idx_assigned_reviewers_unique;
drop index if exists idx_assigned_reviewers_pr_id;
drop index if exists idx_assigned_reviewers_user_id;

drop index if exists idx_prs_merged_at;
drop index if exists idx_prs_created_at;
drop index if exists idx_prs_status;
drop index if exists idx_prs_team_id;
drop index if exists idx_prs_author_id;

drop index if exists idx_team_members_unique;
drop index if exists idx_team_members_team_id;
drop index if exists idx_team_members_user_id;

drop index if exists idx_teams_name;

drop index if exists idx_users_email;
drop index if exists idx_users_is_active;
drop index if exists idx_users_team_name;