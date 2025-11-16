package service

import (
	"context"

	"avitoTest/internal/domain"
	"avitoTest/internal/repository"
)

type UserServiceImpl struct {
	tx   *repository.TxManager
	repo UserRepo
	pr   PRRepo
}

func NewUserService(tx *repository.TxManager, repo UserRepo, pr PRRepo) *UserServiceImpl {
	return &UserServiceImpl{tx: tx, repo: repo, pr: pr}
}

func (s *UserServiceImpl) SetIsActive(ctx context.Context, userID string, active bool) (domain.User, error) {
	var out domain.User
	err := s.tx.WithinTx(ctx, func(ctx context.Context, q repository.Querier) error {
		u, err := s.repo.SetIsActive(ctx, q, userID, active)
		if err != nil {
			return ErrNotFound("user")
		}
		out = u
		return nil
	})
	return out, err
}

func (s *UserServiceImpl) GetReviewPRs(ctx context.Context, userID string) ([]domain.PullRequest, error) {
	var out []domain.PullRequest
	err := s.tx.WithinTx(ctx, func(ctx context.Context, q repository.Querier) error {
		list, err := s.pr.ListUserReviewPRs(ctx, q, userID)
		if err != nil {
			return err
		}
		out = list
		return nil
	})
	return out, err
}
