package models

type Session struct {
	ID      string `json:"id" db:"id"`
	UserID  string `json:"user_id" db:"user_id"`
	IsValid bool   `json:"is_valid" db:"is_valid"`
}
