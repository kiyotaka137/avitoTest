package postgres

import (
	"context"
	"log/slog"

	"avitoTest/internal/domain"
	"avitoTest/internal/repository"

	"github.com/jackc/pgx/v4/pgxpool"
)

type StatsRepo struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func NewStatsRepo(pool *pgxpool.Pool, logger *slog.Logger) *StatsRepo {
	return &StatsRepo{pool: pool, logger: logger}
}

func (r *StatsRepo) CountAssignmentsByUser(ctx context.Context, q repository.Querier) ([]domain.AssignmentByUser, error) {
	rows, err := q.Query(ctx, `
		SELECT u.user_id, u.username, u.team_name, COALESCE(COUNT(prr.pr_id), 0) AS c
		FROM users u
		LEFT JOIN pr_reviewers prr ON prr.reviewer_id = u.user_id
		GROUP BY u.user_id, u.username, u.team_name
		ORDER BY c DESC, u.user_id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.AssignmentByUser
	for rows.Next() {
		var x domain.AssignmentByUser
		if err := rows.Scan(&x.UserID, &x.Username, &x.TeamName, &x.Count); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}

func (r *StatsRepo) CountAssignmentsByPR(ctx context.Context, q repository.Querier) ([]domain.AssignmentByPR, error) {
	rows, err := q.Query(ctx, `
		SELECT p.pr_id, p.pr_name, p.status, COALESCE(COUNT(prr.reviewer_id), 0) AS c
		FROM pull_requests p
		LEFT JOIN pr_reviewers prr ON prr.pr_id = p.pr_id
		GROUP BY p.pr_id, p.pr_name, p.status
		ORDER BY p.pr_id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.AssignmentByPR
	for rows.Next() {
		var x domain.AssignmentByPR
		var status string
		if err := rows.Scan(&x.PullRequestID, &x.PullRequestName, &status, &x.Count); err != nil {
			return nil, err
		}
		x.Status = domain.PRStatus(status)
		out = append(out, x)
	}
	return out, rows.Err()
}
