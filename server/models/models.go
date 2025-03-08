package models

type User struct {
	ID       string `json:"id" db:"id"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type Session struct {
	ID      string `json:"id" db:"id"`
	UserID  string `json:"user_id" db:"user_id"`
	IsValid bool   `json:"is_valid" db:"is_valid"`
}
