package entity

// TODO: А зачем вообще json теги??

type TeamMember struct {
	IsActive *bool  `json:"is_active"`
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}

type Team struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}
