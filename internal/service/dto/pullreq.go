package servdto

import e "app/internal/entity"

type CreatePRInput struct {
	PullReqID string
	NamePR    string
	AuthorID  string
}

type AssignReviewersInput struct {
	PullReqID    string
	ExcludeUsers []string
	AuthorTeam   string
}

type ReassignReviewerInput struct {
	PullReqID string
	RevID     string
	Force     *bool
}

type ReassignReviewerOutput struct {
	PullReq  e.PullRequest
	NewRevID string
}
