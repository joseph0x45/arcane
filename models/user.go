package models

import "time"

type User struct {
	ID        string    `json:"id" db:"id"`
	GithubID  string    `json:"github_id" db:"github_id"`
	Username  string    `json:"username" db:"username"`
	AvatarURL string    `json:"avatar_url" db:"avatar_url"`
	JoinedAt  time.Time `json:"joined_at" db:"joined_at"`
}
