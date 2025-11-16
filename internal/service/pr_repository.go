package service

import (
	"context"
	"time"

	"avitoTest/internal/domain"
	"avitoTest/internal/repository"
)

type PRRepo interface {
	InsertPR(ctx context.Context, q repository.Querier, pr domain.PullRequest) error
	ExistsByID(ctx context.Context, q repository.Querier, id string) (bool, error)
	GetPR(ctx context.Context, q repository.Querier, id string) (domain.PullRequest, error)
	LockPR(ctx context.Context, q repository.Querier, id string) (domain.PullRequest, error)

	AddReviewer(ctx context.Context, q repository.Querier, prID, reviewerID string) error
	RemoveReviewer(ctx context.Context, q repository.Querier, prID, reviewerID string) error
	IsReviewerAssigned(ctx context.Context, q repository.Querier, prID, reviewerID string) (bool, error)

	SetMerged(ctx context.Context, q repository.Querier, id string, t time.Time) error
	ListUserReviewPRs(ctx context.Context, q repository.Querier, userID string) ([]domain.PullRequest, error)
}
