package service

import "fmt"

const (
	CodeTeamExists  = "TEAM_EXISTS"
	CodePRExists    = "PR_EXISTS"
	CodePRMerged    = "PR_MERGED"
	CodeNotAssigned = "NOT_ASSIGNED"
	CodeNoCandidate = "NO_CANDIDATE"
	CodeNotFound    = "NOT_FOUND"
)

type CodedError struct {
	code string
	msg  string
}

func (e CodedError) Error() string { return e.msg }
func (e CodedError) Code() string  { return e.code }

func ErrTeamExists(name string) error {
	return CodedError{code: CodeTeamExists, msg: fmt.Sprintf("team_name '%s' already exists", name)}
}
func ErrPRExists(id string) error {
	return CodedError{code: CodePRExists, msg: fmt.Sprintf("PR '%s' already exists", id)}
}
func ErrPRMerged() error { return CodedError{code: CodePRMerged, msg: "cannot modify merged PR"} }
func ErrNotAssigned() error {
	return CodedError{code: CodeNotAssigned, msg: "reviewer is not assigned to this PR"}
}
func ErrNoCandidate() error {
	return CodedError{code: CodeNoCandidate, msg: "no active replacement candidate in team"}
}
func ErrNotFound(what string) error {
	return CodedError{code: CodeNotFound, msg: fmt.Sprintf("%s not found", what)}
}
