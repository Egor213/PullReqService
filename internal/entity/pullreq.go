package entity

import "time"

type PRStatus string

const (
	StatusOpen   PRStatus = "OPEN"
	StatusMerged PRStatus = "MERGED"
)

type PullRequest struct {
	CreatedAt         *time.Time `db:"created_at"`
	MergedAt          *time.Time `db:"merged_at"`
	PullReqID         string     `db:"pr_id"`
	NamePR            string     `db:"title"`
	AuthorID          string     `db:"author_id"`
	Status            PRStatus   `db:"status"`
	NeedMoreReviewers bool       `db:"need_more_reviewers"`
	Reviewers         []string   `db:"-"`
}

type PullRequestShort struct {
	PullReqID string   `db:"pr_id"`
	NamePR    string   `db:"title"`
	AuthorID  string   `db:"author_id"`
	Status    PRStatus `db:"status"`
}
