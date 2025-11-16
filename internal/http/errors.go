package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	CodeTeamExists  = "TEAM_EXISTS"
	CodePRExists    = "PR_EXISTS"
	CodePRMerged    = "PR_MERGED"
	CodeNotAssigned = "NOT_ASSIGNED"
	CodeNoCandidate = "NO_CANDIDATE"
	CodeNotFound    = "NOT_FOUND"
)

type CodedError interface {
	error
	Code() string
}

func mapError(err error) (status int, code, message string) {
	var ce CodedError
	if errors.As(err, &ce) {
		switch ce.Code() {
		case CodeNotFound:
			return http.StatusNotFound, CodeNotFound, ce.Error()
		case CodeTeamExists:
			return http.StatusBadRequest, CodeTeamExists, ce.Error()
		case CodePRExists:
			return http.StatusConflict, CodePRExists, ce.Error()
		case CodePRMerged:
			return http.StatusConflict, CodePRMerged, ce.Error()
		case CodeNotAssigned:
			return http.StatusConflict, CodeNotAssigned, ce.Error()
		case CodeNoCandidate:
			return http.StatusConflict, CodeNoCandidate, ce.Error()
		default:
			return http.StatusInternalServerError, "INTERNAL", ce.Error()
		}
	}
	// по умолчанию
	return http.StatusInternalServerError, "INTERNAL", "internal error"
}

func respondError(c *gin.Context, status int, code, message string) {
	c.JSON(status, gin.H{
		"error": gin.H{
			"code":    code,
			"message": message,
		},
	})
}
