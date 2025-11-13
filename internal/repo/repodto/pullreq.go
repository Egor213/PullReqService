package repodto

import "app/internal/entity"

type CreatePRInput struct {
	PullReqID string
	NamePR    string
	AuthorID  string
	Status    entity.PRStatus
}
