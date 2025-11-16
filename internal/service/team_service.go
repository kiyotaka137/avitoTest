package service

import (
	"context"

	"avitoTest/internal/domain"
	"avitoTest/internal/repository"
)

type TeamServiceImpl struct {
	tx   *repository.TxManager
	repo TeamRepo
}

func NewTeamService(tx *repository.TxManager, repo TeamRepo) *TeamServiceImpl {
	return &TeamServiceImpl{tx: tx, repo: repo}
}

func (s *TeamServiceImpl) AddTeam(ctx context.Context, team domain.Team) (domain.Team, error) {
	err := s.tx.WithinTx(ctx, func(ctx context.Context, q repository.Querier) error {
		ok, err := s.repo.InsertTeamOrConflict(ctx, q, team.TeamName)
		if err != nil {
			return err
		}
		if !ok {
			return ErrTeamExists(team.TeamName)
		}
		return s.repo.UpsertMembers(ctx, q, team.TeamName, team.Members)
	})
	if err != nil {
		return domain.Team{}, err
	}
	var out domain.Team
	err = s.tx.WithinTx(ctx, func(ctx context.Context, q repository.Querier) error {
		t, err := s.repo.GetTeamWithMembers(ctx, q, team.TeamName)
		if err != nil {
			return err
		}
		out = t
		return nil
	})
	return out, err
}

func (s *TeamServiceImpl) GetTeam(ctx context.Context, teamName string) (domain.Team, error) {
	var out domain.Team
	err := s.tx.WithinTx(ctx, func(ctx context.Context, q repository.Querier) error {
		t, err := s.repo.GetTeamWithMembers(ctx, q, teamName)
		if err != nil {
			return ErrNotFound("team")
		}
		out = t
		return nil
	})
	return out, err
}
