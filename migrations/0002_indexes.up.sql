CREATE INDEX IF NOT EXISTS idx_users_team_active
  ON users (team_name, is_active);

CREATE INDEX IF NOT EXISTS idx_pr_author_status
  ON pull_requests (author_id, status);

CREATE INDEX IF NOT EXISTS idx_pr_reviewers_reviewer
  ON pr_reviewers (reviewer_id);

CREATE INDEX IF NOT EXISTS idx_api_tokens_active_role
  ON api_tokens (is_active, role);
