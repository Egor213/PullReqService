package servdto

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
