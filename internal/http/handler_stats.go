package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"avitoTest/internal/http/dto"
	"avitoTest/internal/ports"
)

type StatsHandler struct {
	Svc ports.StatsService
}

func NewStatsHandler(s ports.StatsService) *StatsHandler { return &StatsHandler{Svc: s} }

func (h *StatsHandler) AssignmentStats(c *gin.Context) {
	stats, err := h.Svc.GetAssignmentStats(c.Request.Context())
	if err != nil {
		status, code, msg := mapError(err)
		respondError(c, status, code, msg)
		return
	}
	out := dto.AssignmentStats{
		ByUser: make([]dto.AssignmentByUser, 0, len(stats.ByUser)),
		ByPR:   make([]dto.AssignmentByPR, 0, len(stats.ByPR)),
	}
	for _, u := range stats.ByUser {
		out.ByUser = append(out.ByUser, dto.AssignmentByUser{
			UserID: u.UserID, Username: u.Username, TeamName: u.TeamName, AssignedCount: u.Count,
		})
	}
	for _, p := range stats.ByPR {
		out.ByPR = append(out.ByPR, dto.AssignmentByPR{
			PullRequestID: p.PullRequestID, PullRequestName: p.PullRequestName,
			Status: string(p.Status), AssignedCount: p.Count,
		})
	}
	c.JSON(http.StatusOK, out)
}
