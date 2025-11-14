package httpdto

type ReviewerStatsDTO struct {
	UserID      string `json:"user_id"`
	Assignments int    `json:"assignments"`
}

type PRStatsDTO struct {
	PullReqID   string `json:"pr_id"`
	Assignments int    `json:"assignments"`
}

type GetStatsOutput struct {
	ByUsers []ReviewerStatsDTO `json:"by_users"`
	ByPRs   []PRStatsDTO       `json:"by_prs"`
}
