package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"avitoTest/internal/http/dto"
	"avitoTest/internal/ports"
)

type PullRequestHandler struct {
	Svc ports.PRService
}

func NewPullRequestHandler(s ports.PRService) *PullRequestHandler { return &PullRequestHandler{Svc: s} }

type createReq struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}
type mergeReq struct {
	PullRequestID string `json:"pull_request_id"`
}
type reassignReq struct {
	PullRequestID string `json:"pull_request_id"`
	OldUserID     string `json:"old_user_id"`
}

func (h *PullRequestHandler) Create(c *gin.Context) {
	var in createReq
	if err := c.ShouldBindJSON(&in); err != nil || in.PullRequestID == "" || in.PullRequestName == "" || in.AuthorID == "" {
		respondError(c, http.StatusBadRequest, "BAD_REQUEST", "pull_request_id, pull_request_name, author_id are required")
		return
	}
	pr, err := h.Svc.CreatePR(c.Request.Context(), ports.CreatePRInput{
		PullRequestID:   in.PullRequestID,
		PullRequestName: in.PullRequestName,
		AuthorID:        in.AuthorID,
	})
	if err != nil {
		status, code, msg := mapError(err)
		respondError(c, status, code, msg)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"pr": dto.FromDomainPR(pr)})
}

func (h *PullRequestHandler) Merge(c *gin.Context) {
	var in mergeReq
	if err := c.ShouldBindJSON(&in); err != nil || in.PullRequestID == "" {
		respondError(c, http.StatusBadRequest, "BAD_REQUEST", "pull_request_id is required")
		return
	}
	pr, err := h.Svc.MergePR(c.Request.Context(), in.PullRequestID)
	if err != nil {
		status, code, msg := mapError(err)
		respondError(c, status, code, msg)
		return
	}
	c.JSON(http.StatusOK, gin.H{"pr": dto.FromDomainPR(pr)})
}

func (h *PullRequestHandler) Reassign(c *gin.Context) {
	var in reassignReq
	if err := c.ShouldBindJSON(&in); err != nil || in.PullRequestID == "" || in.OldUserID == "" {
		respondError(c, http.StatusBadRequest, "BAD_REQUEST", "pull_request_id and old_user_id are required")
		return
	}
	pr, replacedBy, err := h.Svc.ReassignReviewer(c.Request.Context(), in.PullRequestID, in.OldUserID)
	if err != nil {
		status, code, msg := mapError(err)
		respondError(c, status, code, msg)
		return
	}
	c.JSON(http.StatusOK, gin.H{"pr": dto.FromDomainPR(pr), "replaced_by": replacedBy})
}
