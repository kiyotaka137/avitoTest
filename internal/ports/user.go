package ports

import (
	"context"

	"avitoTest/internal/domain"
)

type UserService interface {
	SetIsActive(ctx context.Context, userID string, active bool) (domain.User, error)
	GetReviewPRs(ctx context.Context, userID string) ([]domain.PullRequest, error)
}
