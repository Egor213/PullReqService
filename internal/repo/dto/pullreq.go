package repodto

import e "app/internal/entity"

type CreatePRInput struct {
	PullReqID string
	NamePR    string
	AuthorID  string
	Status    e.PRStatus
}

type ChangeReviewerInput struct {
	PullReqID   string
	NewReviewer string
	OldReviewer string
}

type GetAutAndStOutput struct {
	AuthorID string
	Status   e.PRStatus
}
