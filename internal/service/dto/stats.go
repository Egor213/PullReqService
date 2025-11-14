package servdto

type ReviewerStatsDTO struct {
	UserID      string
	Assignments int
}

type PRStatsDTO struct {
	PullReqID   string
	Assignments int
}

type GetStatsOutput struct {
	ByUsers []ReviewerStatsDTO
	ByPRs   []PRStatsDTO
}
