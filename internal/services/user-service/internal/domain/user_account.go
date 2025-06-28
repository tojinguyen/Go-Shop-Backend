package domain

type UserAccount struct {
	Id             string `json:"id"`
	Email          string `json:"email"`
	HashedPassword string `json:"-"` // Don't serialize password in JSON
	LastLoginAt    string `json:"last_login_at"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}
