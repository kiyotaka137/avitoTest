DROP TABLE IF EXISTS api_tokens;
DROP TABLE IF EXISTS pr_reviewers;
DROP TABLE IF EXISTS pull_requests;
DROP TRIGGER IF EXISTS trg_users_updated_at ON users;
DROP FUNCTION IF EXISTS set_updated_at;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS teams;

DO $$
BEGIN
  IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'pr_status') THEN
    DROP TYPE pr_status;
  END IF;
END$$;
