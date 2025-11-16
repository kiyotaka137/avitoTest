package dto

type AssignmentByUser struct {
	UserID        string `json:"user_id"`
	Username      string `json:"username"`
	TeamName      string `json:"team_name"`
	AssignedCount int    `json:"assigned_count"`
}

type AssignmentByPR struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	Status          string `json:"status"`
	AssignedCount   int    `json:"assigned_count"`
}

type AssignmentStats struct {
	ByUser []AssignmentByUser `json:"by_user"`
	ByPR   []AssignmentByPR   `json:"by_pr"`
}
