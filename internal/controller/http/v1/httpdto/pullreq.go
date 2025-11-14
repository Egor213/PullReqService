package httpdto

type CreatePRInput struct {
	PullReqID string `json:"pull_request_id" validate:"required"`
	NamePR    string `json:"pull_request_name" validate:"required"`
	AuthorID  string `json:"author_id" validate:"required"`
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
	PRID string `query:"pr_id" validate:"required"`
}

type GetPROutput struct {
	PullRequestDTO
	NeedMoreReviewers *bool `json:"need_more_reviewers"`
}

type ReassignReviewerInput struct {
	PullReqID   string `json:"pull_request_id"`
	OldReviewer string `json:"old_reviewer_id"`
}

type ReassignReviewerOutput struct {
	PullReq     PullRequestDTO `json:"pr"`
	NewReviewer string         `json:"replaced_by"`
}
