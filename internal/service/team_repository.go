package service

import (
	"context"

	"avitoTest/internal/domain"
	"avitoTest/internal/repository"
)

type TeamRepo interface {
	InsertTeamOrConflict(ctx context.Context, q repository.Querier, teamName string) (bool, error)
	GetTeamWithMembers(ctx context.Context, q repository.Querier, teamName string) (domain.Team, error)
	UpsertMembers(ctx context.Context, q repository.Querier, teamName string, members []domain.TeamMember) error
}
