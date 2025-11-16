package postgres

import (
	"context"
	"log/slog"
	"time"

	"avitoTest/internal/domain"
	"avitoTest/internal/repository"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PRRepo struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func NewPRRepo(pool *pgxpool.Pool, logger *slog.Logger) *PRRepo {
	return &PRRepo{pool: pool, logger: logger}
}

func (r *PRRepo) InsertPR(ctx context.Context, q repository.Querier, pr domain.PullRequest) error {
	_, err := q.Exec(ctx, `
		INSERT INTO pull_requests (pr_id, pr_name, author_id, status, created_at, merged_at)
		VALUES ($1, $2, $3, $4, NOW(), NULL)
	`, pr.PullRequestID, pr.PullRequestName, pr.AuthorID, string(pr.Status))
	return err
}

func (r *PRRepo) ExistsByID(ctx context.Context, q repository.Querier, id string) (bool, error) {
	var ex bool
	err := q.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM pull_requests WHERE pr_id=$1)`, id).Scan(&ex)
	return ex, err
}

func (r *PRRepo) GetPR(ctx context.Context, q repository.Querier, id string) (domain.PullRequest, error) {
	var pr domain.PullRequest
	var status string
	var created time.Time
	var merged *time.Time

	if err := q.QueryRow(ctx, `
		SELECT pr_id, pr_name, author_id, status, created_at, merged_at
		FROM pull_requests WHERE pr_id=$1
	`, id).Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &status, &created, &merged); err != nil {
		return domain.PullRequest{}, ErrNotFound("pull_request")
	}
	pr.Status = domain.PRStatus(status)
	pr.CreatedAt = &created
	pr.MergedAt = merged

	rows, err := q.Query(ctx, `SELECT reviewer_id FROM pr_reviewers WHERE pr_id=$1 ORDER BY reviewer_id ASC`, id)
	if err != nil {
		return domain.PullRequest{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var rid string
		if err := rows.Scan(&rid); err != nil {
			return domain.PullRequest{}, err
		}
		pr.AssignedReviewers = append(pr.AssignedReviewers, rid)
	}
	return pr, rows.Err()
}

func (r *PRRepo) AddReviewer(ctx context.Context, q repository.Querier, prID, reviewerID string) error {
	_, err := q.Exec(ctx, `
		INSERT INTO pr_reviewers (pr_id, reviewer_id)
		VALUES ($1, $2) ON CONFLICT DO NOTHING
	`, prID, reviewerID)
	return err
}

func (r *PRRepo) RemoveReviewer(ctx context.Context, q repository.Querier, prID, reviewerID string) error {
	_, err := q.Exec(ctx, `DELETE FROM pr_reviewers WHERE pr_id=$1 AND reviewer_id=$2`, prID, reviewerID)
	return err
}

func (r *PRRepo) IsReviewerAssigned(ctx context.Context, q repository.Querier, prID, reviewerID string) (bool, error) {
	var ex bool
	err := q.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM pr_reviewers WHERE pr_id=$1 AND reviewer_id=$2)
	`, prID, reviewerID).Scan(&ex)
	return ex, err
}

func (r *PRRepo) LockPR(ctx context.Context, q repository.Querier, id string) (domain.PullRequest, error) {
	var pr domain.PullRequest
	var status string
	var created time.Time
	var merged *time.Time

	if err := q.QueryRow(ctx, `
		SELECT pr_id, pr_name, author_id, status, created_at, merged_at
		FROM pull_requests WHERE pr_id=$1
		FOR UPDATE
	`, id).Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &status, &created, &merged); err != nil {
		return domain.PullRequest{}, ErrNotFound("pull_request")
	}
	pr.Status = domain.PRStatus(status)
	pr.CreatedAt = &created
	pr.MergedAt = merged

	rows, err := q.Query(ctx, `SELECT reviewer_id FROM pr_reviewers WHERE pr_id=$1 ORDER BY reviewer_id ASC`, id)
	if err != nil {
		return domain.PullRequest{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var rid string
		if err := rows.Scan(&rid); err != nil {
			return domain.PullRequest{}, err
		}
		pr.AssignedReviewers = append(pr.AssignedReviewers, rid)
	}
	return pr, rows.Err()
}

func (r *PRRepo) SetMerged(ctx context.Context, q repository.Querier, id string, t time.Time) error {
	_, err := q.Exec(ctx, `
		UPDATE pull_requests
		SET status='MERGED', merged_at=$2
		WHERE pr_id=$1
	`, id, t.UTC())
	return err
}

func (r *PRRepo) ListUserReviewPRs(ctx context.Context, q repository.Querier, userID string) ([]domain.PullRequest, error) {
	rows, err := q.Query(ctx, `
		SELECT p.pr_id, p.pr_name, p.author_id, p.status, p.created_at, p.merged_at
		FROM pr_reviewers r
		JOIN pull_requests p ON p.pr_id = r.pr_id
		WHERE r.reviewer_id=$1
		ORDER BY p.created_at DESC, p.pr_id ASC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.PullRequest
	for rows.Next() {
		var pr domain.PullRequest
		var status string
		var created time.Time
		var merged *time.Time
		if err := rows.Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &status, &created, &merged); err != nil {
			return nil, err
		}
		pr.Status = domain.PRStatus(status)
		pr.CreatedAt = &created
		pr.MergedAt = merged
		out = append(out, pr)
	}
	return out, rows.Err()
}
