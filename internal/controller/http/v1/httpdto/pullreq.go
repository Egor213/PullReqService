package httpdto

import "time"

type CreatePRInput struct {
	PullReqID string `json:"pull_request_id" validate:"required,max=100"`
	NamePR    string `json:"pull_request_name" validate:"required,max=100"`
	AuthorID  string `json:"author_id" validate:"required,max=100"`
}

type CreatePROutput struct {
	PullReq PullRequestDTO `json:"pr"`
}

type PullRequestDTO struct {
	PullReqID string   `json:"pull_request_id"`
	NamePR    string   `json:"pull_request_name"`
	AuthorID  string   `json:"author_id"`
	Status    string   `json:"status"`
	Reviewers []string `json:"assigned_reviewers"`
}

type GetPRInput struct {
	PRID string `query:"pr_id" validate:"required,max=100"`
}

type GetPROutput struct {
	PullRequestDTO
	NeedMoreReviewers *bool      `json:"need_more_reviewers"`
	MergedAt          *time.Time `json:"mergedAt"`
}

type ReassignReviewerInput struct {
	PullReqID   string `json:"pull_request_id" validate:"required,max=100"`
	OldReviewer string `json:"old_reviewer_id" validate:"required,max=100"`
}

type ReassignReviewerOutput struct {
	PullReq     PullRequestDTO `json:"pr"`
	NewReviewer string         `json:"replaced_by"`
}

type PullRequestShortDTO struct {
	PullReqID string `json:"pull_request_id"`
	NamePR    string `json:"pull_request_name"`
	AuthorID  string `json:"author_id"`
	Status    string `json:"status"`
}

type GetReviewInput struct {
	UserID string `query:"user_id" validate:"required,max=100"`
}

type GetReviewOutput struct {
	UserID  string                `json:"user_id"`
	PullReq []PullRequestShortDTO `json:"pull_requests"`
}

type MergePRInput struct {
	PullReqID string `json:"pull_request_id" validate:"required,max=100"`
}

type MergePRODTO struct {
	PullRequestDTO
	MergedAt *time.Time `json:"mergedAt"`
}

type MergePROutput struct {
	PullReq MergePRODTO `json:"pr"`
}
