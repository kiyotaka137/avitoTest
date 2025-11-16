package ports

import (
	"avitoTest/internal/domain"
	"context"
)

type StatsService interface {
	GetAssignmentStats(ctx context.Context) (domain.AssignmentStats, error)
}
