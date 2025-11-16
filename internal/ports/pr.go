package ports

import (
	"context"

	"avitoTest/internal/domain"
)

type CreatePRInput struct {
	PullRequestID   string
	PullRequestName string
	AuthorID        string
}

type PRService interface {
	CreatePR(ctx context.Context, in CreatePRInput) (domain.PullRequest, error)
	MergePR(ctx context.Context, prID string) (domain.PullRequest, error)
	ReassignReviewer(ctx context.Context, prID, oldUserID string) (domain.PullRequest, string, error)
}
