package dto

import (
	"avitoTest/internal/domain"
	"time"
)

func FromDomainUser(u domain.User) User {
	return User{
		UserID:   u.UserID,
		Username: u.Username,
		TeamName: u.TeamName,
		IsActive: u.IsActive,
	}
}

func FromDomainTeam(t domain.Team) Team {
	ms := make([]TeamMember, 0, len(t.Members))
	for _, m := range t.Members {
		ms = append(ms, TeamMember{UserID: m.UserID, Username: m.Username, IsActive: m.IsActive})
	}
	return Team{TeamName: t.TeamName, Members: ms}
}

func ToDomainTeam(t Team) domain.Team {
	ms := make([]domain.TeamMember, 0, len(t.Members))
	for _, m := range t.Members {
		ms = append(ms, domain.TeamMember{UserID: m.UserID, Username: m.Username, IsActive: m.IsActive})
	}
	return domain.Team{TeamName: t.TeamName, Members: ms}
}

func toISO(tt *time.Time) *string {
	if tt == nil {
		return nil
	}
	s := tt.UTC().Format(time.RFC3339)
	return &s
}

func FromDomainPR(pr domain.PullRequest) PullRequest {
	return PullRequest{
		PullRequestID:     pr.PullRequestID,
		PullRequestName:   pr.PullRequestName,
		AuthorID:          pr.AuthorID,
		Status:            string(pr.Status),
		AssignedReviewers: append([]string(nil), pr.AssignedReviewers...),
		CreatedAt:         toISO(pr.CreatedAt),
		MergedAt:          toISO(pr.MergedAt),
	}
}

func ShortFromDomainPR(pr domain.PullRequest) PullRequestShort {
	return PullRequestShort{
		PullRequestID:   pr.PullRequestID,
		PullRequestName: pr.PullRequestName,
		AuthorID:        pr.AuthorID,
		Status:          string(pr.Status),
	}
}
