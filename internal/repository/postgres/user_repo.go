package postgres

import (
	"context"
	"log/slog"

	"avitoTest/internal/domain"
	"avitoTest/internal/repository"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserRepo struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func NewUserRepo(pool *pgxpool.Pool, logger *slog.Logger) *UserRepo {
	return &UserRepo{pool: pool, logger: logger}
}

func (r *UserRepo) SetIsActive(ctx context.Context, q repository.Querier, userID string, active bool) (domain.User, error) {
	_, err := q.Exec(ctx, `UPDATE users SET is_active=$2 WHERE user_id=$1`, userID, active)
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = q.QueryRow(ctx, `
		SELECT user_id, username, team_name, is_active
		FROM users WHERE user_id=$1
	`, userID).Scan(&u.UserID, &u.Username, &u.TeamName, &u.IsActive)
	if err != nil {
		return domain.User{}, ErrNotFound("user")
	}
	return u, nil
}

func (r *UserRepo) GetByID(ctx context.Context, q repository.Querier, id string) (domain.User, error) {
	var u domain.User
	err := q.QueryRow(ctx, `
		SELECT user_id, username, team_name, is_active
		FROM users WHERE user_id=$1
	`, id).Scan(&u.UserID, &u.Username, &u.TeamName, &u.IsActive)
	if err != nil {
		return domain.User{}, ErrNotFound("user")
	}
	return u, nil
}

func (r *UserRepo) ListActiveByTeamExcept(ctx context.Context, q repository.Querier, teamName string, excludeIDs []string, limit int) ([]domain.User, error) {
	var arr pgtype.TextArray
	if err := arr.Set(excludeIDs); err != nil {
		return nil, err
	}

	rows, err := q.Query(ctx, `
		SELECT user_id, username, team_name, is_active
		FROM users
		WHERE team_name=$1
		  AND is_active=TRUE
		  AND NOT (user_id = ANY($2))
		ORDER BY user_id ASC
		LIMIT $3
	`, teamName, arr, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.UserID, &u.Username, &u.TeamName, &u.IsActive); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}
