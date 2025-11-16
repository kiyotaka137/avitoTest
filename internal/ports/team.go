package ports

import (
	"context"

	"avitoTest/internal/domain"
)

type TeamService interface {
	AddTeam(ctx context.Context, team domain.Team) (domain.Team, error)
	GetTeam(ctx context.Context, teamName string) (domain.Team, error)
}
