package postgres

import (
	"context"
	"log/slog"

	"avitoTest/internal/domain"
	"avitoTest/internal/repository"

	"github.com/jackc/pgx/v4/pgxpool"
)

type TeamRepo struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func NewTeamRepo(pool *pgxpool.Pool, logger *slog.Logger) *TeamRepo {
	return &TeamRepo{pool: pool, logger: logger}
}

func (r *TeamRepo) InsertTeamOrConflict(ctx context.Context, q repository.Querier, teamName string) (bool, error) {
	tag, err := q.Exec(ctx, `INSERT INTO teams (team_name) VALUES ($1) ON CONFLICT DO NOTHING`, teamName)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() == 1, nil
}

func (r *TeamRepo) GetTeamWithMembers(ctx context.Context, q repository.Querier, teamName string) (domain.Team, error) {
	rows, err := q.Query(ctx, `
		SELECT u.user_id, u.username, u.is_active
		FROM users u
		WHERE u.team_name = $1
		ORDER BY u.user_id ASC
	`, teamName)
	if err != nil {
		return domain.Team{}, err
	}
	defer rows.Close()

	var members []domain.TeamMember
	for rows.Next() {
		var m domain.TeamMember
		if err := rows.Scan(&m.UserID, &m.Username, &m.IsActive); err != nil {
			return domain.Team{}, err
		}
		members = append(members, m)
	}
	if rows.Err() != nil {
		return domain.Team{}, rows.Err()
	}
	var exists bool
	if err := q.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM teams WHERE team_name=$1)`, teamName).Scan(&exists); err != nil {
		return domain.Team{}, err
	}
	if !exists {
		return domain.Team{}, ErrNotFound("team")
	}
	return domain.Team{TeamName: teamName, Members: members}, nil
}

func (r *TeamRepo) UpsertMembers(ctx context.Context, q repository.Querier, teamName string, members []domain.TeamMember) error {
	for _, m := range members {
		_, err := q.Exec(ctx, `
			INSERT INTO users (user_id, username, team_name, is_active)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (user_id) DO UPDATE
			SET username = EXCLUDED.username,
			    team_name = EXCLUDED.team_name,
			    is_active = EXCLUDED.is_active
		`, m.UserID, m.Username, teamName, m.IsActive)
		if err != nil {
			return err
		}
	}
	return nil
}
