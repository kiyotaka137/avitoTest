package service

import (
	"context"

	"avitoTest/internal/domain"
	"avitoTest/internal/repository"
)

type UserRepo interface {
	SetIsActive(ctx context.Context, q repository.Querier, userID string, active bool) (domain.User, error)
	GetByID(ctx context.Context, q repository.Querier, id string) (domain.User, error)
	ListActiveByTeamExcept(ctx context.Context, q repository.Querier, teamName string, excludeIDs []string, limit int) ([]domain.User, error)
}
