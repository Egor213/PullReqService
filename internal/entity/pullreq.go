package entity

import "time"

type PRStatus string

const (
	StatusOpen   PRStatus = "OPEN"
	StatusMerged PRStatus = "MERGED"
)

type PullRequest struct {
	CreatedAt         *time.Time
	MergedAt          *time.Time
	PullReqID         string
	NamePR            string
	AuthorID          string
	Status            PRStatus
	NeedMoreReviewers bool
	Reviewers         []string
}
