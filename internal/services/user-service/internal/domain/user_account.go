package domain

type UserAccount struct {
	Id          string `json:"id"`
	Email       string `json:"email"`
	LastLoginAt string `json:"last_login_at"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
