package repodto

type ReviewerStatsOutput struct {
	UserID      string `db:"user_id"`
	Assignments int    `db:"assignments"`
}

type PRStatsOutput struct {
	PullReqID   string `db:"pr_id"`
	Assignments int    `db:"assignments"`
}
