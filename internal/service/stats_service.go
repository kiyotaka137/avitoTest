package service

import (
	"context"

	"avitoTest/internal/domain"
	"avitoTest/internal/repository"
)

type StatsRepo interface {
	CountAssignmentsByUser(ctx context.Context, q repository.Querier) ([]domain.AssignmentByUser, error)
	CountAssignmentsByPR(ctx context.Context, q repository.Querier) ([]domain.AssignmentByPR, error)
}

type StatsServiceImpl struct {
	tx   *repository.TxManager
	repo StatsRepo
}

func NewStatsService(tx *repository.TxManager, repo StatsRepo) *StatsServiceImpl {
	return &StatsServiceImpl{tx: tx, repo: repo}
}

func (s *StatsServiceImpl) GetAssignmentStats(ctx context.Context) (domain.AssignmentStats, error) {
	var out domain.AssignmentStats
	err := s.tx.WithinTx(ctx, func(ctx context.Context, q repository.Querier) error {
		users, err := s.repo.CountAssignmentsByUser(ctx, q)
		if err != nil {
			return err
		}
		prs, err := s.repo.CountAssignmentsByPR(ctx, q)
		if err != nil {
			return err
		}
		out.ByUser = users
		out.ByPR = prs
		return nil
	})
	return out, err
}
