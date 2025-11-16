package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"avitoTest/internal/http/dto"
	"avitoTest/internal/ports"
)

type TeamHandler struct {
	Svc ports.TeamService
}

func NewTeamHandler(s ports.TeamService) *TeamHandler { return &TeamHandler{Svc: s} }

func (h *TeamHandler) AddTeam(c *gin.Context) {
	var in dto.Team
	if err := c.ShouldBindJSON(&in); err != nil || in.TeamName == "" {
		respondError(c, http.StatusBadRequest, "BAD_REQUEST", "invalid body")
		return
	}
	for _, m := range in.Members {
		if m.UserID == "" || m.Username == "" {
			respondError(c, http.StatusBadRequest, "BAD_REQUEST", "member.user_id and member.username are required")
			return
		}
	}
	created, err := h.Svc.AddTeam(c.Request.Context(), dto.ToDomainTeam(in))
	if err != nil {
		status, code, msg := mapError(err)
		respondError(c, status, code, msg)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"team": dto.FromDomainTeam(created)})
}

func (h *TeamHandler) GetTeam(c *gin.Context) {
	name := c.Query("team_name")
	if name == "" {
		respondError(c, http.StatusBadRequest, "BAD_REQUEST", "team_name is required")
		return
	}
	team, err := h.Svc.GetTeam(c.Request.Context(), name)
	if err != nil {
		status, code, msg := mapError(err)
		respondError(c, status, code, msg)
		return
	}
	c.JSON(http.StatusOK, dto.FromDomainTeam(team))
}
