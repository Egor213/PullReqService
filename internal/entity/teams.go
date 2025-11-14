package entity

type TeamMember struct {
	IsActive *bool  `db:"is_active"`
	UserID   string `db:"user_id"`
	Username string `db:"username"`
}

type Team struct {
	TeamName string       `db:"team_name"`
	Members  []TeamMember `db:"-"`
}
