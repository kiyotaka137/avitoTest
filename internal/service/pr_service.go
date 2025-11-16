package service

import (
	"context"
	"time"

	"avitoTest/internal/domain"
	"avitoTest/internal/ports"
	"avitoTest/internal/repository"
)

type PRServiceImpl struct {
	tx     *repository.TxManager
	prRepo PRRepo
	user   UserRepo
	team   TeamRepo
}

func NewPRService(tx *repository.TxManager, pr PRRepo, user UserRepo, team TeamRepo) *PRServiceImpl {
	return &PRServiceImpl{tx: tx, prRepo: pr, user: user, team: team}
}

func (s *PRServiceImpl) CreatePR(ctx context.Context, in ports.CreatePRInput) (domain.PullRequest, error) {
	var out domain.PullRequest
	err := s.tx.WithinTx(ctx, func(ctx context.Context, q repository.Querier) error {
		author, err := s.user.GetByID(ctx, q, in.AuthorID)
		if err != nil {
			return ErrNotFound("author")
		}
		if _, err := s.repoTeamExists(ctx, q, author.TeamName); err != nil {
			return err
		}
		exists, err := s.prRepo.ExistsByID(ctx, q, in.PullRequestID)
		if err != nil {
			return err
		}
		if exists {
			return ErrPRExists(in.PullRequestID)
		}
		pr := domain.PullRequest{
			PullRequestID:   in.PullRequestID,
			PullRequestName: in.PullRequestName,
			AuthorID:        in.AuthorID,
			Status:          domain.PRStatusOpen,
		}
		if err := s.prRepo.InsertPR(ctx, q, pr); err != nil {
			return err
		}
		exclude := []string{author.UserID}
		cands, err := s.user.ListActiveByTeamExcept(ctx, q, author.TeamName, exclude, 2)
		if err != nil {
			return err
		}
		for _, c := range cands {
			if err := s.prRepo.AddReviewer(ctx, q, in.PullRequestID, c.UserID); err != nil {
				return err
			}
			pr.AssignedReviewers = append(pr.AssignedReviewers, c.UserID)
		}
		out = pr
		return nil
	})
	if err != nil {
		return domain.PullRequest{}, err
	}
	_ = s.tx.WithinTx(ctx, func(ctx context.Context, q repository.Querier) error {
		p, err := s.prRepo.GetPR(ctx, q, in.PullRequestID)
		if err == nil {
			out = p
		}
		return nil
	})
	return out, nil
}

func (s *PRServiceImpl) MergePR(ctx context.Context, prID string) (domain.PullRequest, error) {
	var out domain.PullRequest
	err := s.tx.WithinTx(ctx, func(ctx context.Context, q repository.Querier) error {
		pr, err := s.prRepo.LockPR(ctx, q, prID)
		if err != nil {
			return ErrNotFound("pull_request")
		}
		if pr.Status == domain.PRStatusMerged {
			out = pr
			return nil
		}
		now := time.Now().UTC()
		if err := s.prRepo.SetMerged(ctx, q, prID, now); err != nil {
			return err
		}
		merged, err := s.prRepo.GetPR(ctx, q, prID)
		if err != nil {
			return err
		}
		out = merged
		return nil
	})
	return out, err
}

func (s *PRServiceImpl) ReassignReviewer(ctx context.Context, prID, oldUserID string) (domain.PullRequest, string, error) {
	var out domain.PullRequest
	var replacedBy string
	err := s.tx.WithinTx(ctx, func(ctx context.Context, q repository.Querier) error {
		pr, err := s.prRepo.LockPR(ctx, q, prID)
		if err != nil {
			return ErrNotFound("pull_request")
		}
		if pr.Status == domain.PRStatusMerged {
			return ErrPRMerged()
		}
		assigned, err := s.prRepo.IsReviewerAssigned(ctx, q, prID, oldUserID)
		if err != nil {
			return err
		}
		if !assigned {
			return ErrNotAssigned()
		}
		author, err := s.user.GetByID(ctx, q, pr.AuthorID)
		if err != nil {
			return ErrNotFound("author")
		}
		exclude := append([]string{author.UserID, oldUserID}, pr.AssignedReviewers...)
		exMap := map[string]struct{}{}
		uniq := make([]string, 0, len(exclude))
		for _, id := range exclude {
			if _, ok := exMap[id]; !ok {
				exMap[id] = struct{}{}
				uniq = append(uniq, id)
			}
		}
		cands, err := s.user.ListActiveByTeamExcept(ctx, q, author.TeamName, uniq, 1)
		if err != nil {
			return err
		}
		if len(cands) == 0 {
			return ErrNoCandidate()
		}
		newID := cands[0].UserID
		if err := s.prRepo.RemoveReviewer(ctx, q, prID, oldUserID); err != nil {
			return err
		}
		if err := s.prRepo.AddReviewer(ctx, q, prID, newID); err != nil {
			return err
		}
		replacedBy = newID
		updated, err := s.prRepo.GetPR(ctx, q, prID)
		if err != nil {
			return err
		}
		out = updated
		return nil
	})
	return out, replacedBy, err
}

func (s *PRServiceImpl) repoTeamExists(ctx context.Context, q repository.Querier, teamName string) (bool, error) {
	_, err := s.team.GetTeamWithMembers(ctx, q, teamName)
	if err != nil {
		return false, ErrNotFound("team")
	}
	return true, nil
}
