package domain

type AssignmentByUser struct {
	UserID   string
	Username string
	TeamName string
	Count    int
}

type AssignmentByPR struct {
	PullRequestID   string
	PullRequestName string
	Status          PRStatus
	Count           int
}

type AssignmentStats struct {
	ByUser []AssignmentByUser
	ByPR   []AssignmentByPR
}
