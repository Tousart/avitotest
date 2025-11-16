package models

// ErrorResponse

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Team

type TeamMember struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type Team struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}

// User

type User struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

// PR

type PullRequest struct {
	PullRequestID     string   `json:"pull_request_id"`
	PullRequestName   string   `json:"pull_request_name"`
	AuthorID          string   `json:"author_id"`
	Status            string   `json:"status"`
	AssignedReviewers []string `json:"assigned_reviewers"`
	CreatedAt         string   `json:"created_at"`
	MergedAt          string   `json:"merged_at"`
}

type PullRequestShort struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}

// Активность - добавил от себя
// Пользователь, сколько пулл реквестов ревьюит, сколько из них MERGED и сколько из них OPEN
type UserActivity struct {
	UserID       string `json:"user_id"`
	Username     string `json:"username"`
	PullRequests int    `json:"pull_requests"`
	MergedPR     int    `json:"merged_pr"`
	OpenPR       int    `json:"open_pr"`
}
