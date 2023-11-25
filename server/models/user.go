package models

type User struct {
	Id        string `json:"id" db:"id"`
	GithubId  string `json:"github_id" db:"github_id"`
	Email     string `json:"email" db:"email"`
	Username  string `json:"username" db:"username"`
	AvatarURL string `json:"avatar_url" db:"avatar_url"`
}
