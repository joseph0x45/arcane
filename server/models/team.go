package models

type Team struct {
	Id    string `json:"id" db:"id"`
	Name  string `json:"name" db:"id"`
	Owner string `json:"owner" db:"owner"`
	Plan  string `json:"plan" db:"plan"`
}
