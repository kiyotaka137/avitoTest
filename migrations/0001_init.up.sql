
DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'pr_status') THEN
    CREATE TYPE pr_status AS ENUM ('OPEN', 'MERGED');
  END IF;
END$$;

CREATE TABLE IF NOT EXISTS teams (
  team_name TEXT PRIMARY KEY,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS users (
  user_id   TEXT PRIMARY KEY,
  username  TEXT NOT NULL,
  team_name TEXT NOT NULL REFERENCES teams(team_name) ON UPDATE CASCADE ON DELETE RESTRICT,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_users_updated_at ON users;
CREATE TRIGGER trg_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE IF NOT EXISTS pull_requests (
  pr_id       TEXT PRIMARY KEY,
  pr_name     TEXT NOT NULL,
  author_id   TEXT NOT NULL REFERENCES users(user_id) ON UPDATE CASCADE ON DELETE RESTRICT,
  status      pr_status NOT NULL DEFAULT 'OPEN',
  created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  merged_at   TIMESTAMPTZ NULL
);

CREATE TABLE IF NOT EXISTS pr_reviewers (
  pr_id       TEXT NOT NULL REFERENCES pull_requests(pr_id) ON UPDATE CASCADE ON DELETE CASCADE,
  reviewer_id TEXT NOT NULL REFERENCES users(user_id)       ON UPDATE CASCADE ON DELETE RESTRICT,
  PRIMARY KEY (pr_id, reviewer_id)
);

CREATE TABLE IF NOT EXISTS api_tokens (
  token      TEXT PRIMARY KEY,
  user_id    TEXT NULL REFERENCES users(user_id) ON UPDATE CASCADE ON DELETE SET NULL,
  role       TEXT NOT NULL CHECK (role IN ('admin','user')),
  is_active  BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
