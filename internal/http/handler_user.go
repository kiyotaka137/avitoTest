package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"avitoTest/internal/http/dto"
	"avitoTest/internal/ports"
)

type UserHandler struct {
	Svc ports.UserService
}

func NewUserHandler(s ports.UserService) *UserHandler { return &UserHandler{Svc: s} }

type setIsActiveReq struct {
	UserID   string `json:"user_id"`
	IsActive *bool  `json:"is_active"`
}

func (h *UserHandler) SetIsActive(c *gin.Context) {
	var in setIsActiveReq
	if err := c.ShouldBindJSON(&in); err != nil || in.UserID == "" || in.IsActive == nil {
		respondError(c, http.StatusBadRequest, "BAD_REQUEST", "user_id and is_active are required")
		return
	}
	u, err := h.Svc.SetIsActive(c.Request.Context(), in.UserID, *in.IsActive)
	if err != nil {
		status, code, msg := mapError(err)
		respondError(c, status, code, msg)
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": dto.FromDomainUser(u)})
}

func (h *UserHandler) GetReviewPRs(c *gin.Context) {
	uid := c.Query("user_id")
	if uid == "" {
		respondError(c, http.StatusBadRequest, "BAD_REQUEST", "user_id is required")
		return
	}
	prs, err := h.Svc.GetReviewPRs(c.Request.Context(), uid)
	if err != nil {
		status, code, msg := mapError(err)
		respondError(c, status, code, msg)
		return
	}
	short := make([]dto.PullRequestShort, 0, len(prs))
	for _, p := range prs {
		short = append(short, dto.ShortFromDomainPR(p))
	}
	c.JSON(http.StatusOK, gin.H{"user_id": uid, "pull_requests": short})
}
