package store

import (
	"server/models"

	"github.com/jmoiron/sqlx"
)

type Users struct {
	db *sqlx.DB
}

func NewUsersStore(db *sqlx.DB) *Users {
	return &Users{
		db: db,
	}
}

func (s *Users) Insert(u *models.User) error {
	_, err := s.db.NamedExec(
		"insert into users(id, email, github_id, username, avatar_url) values(:id, :email, :github_id, :username, :avatar_url)",
		u,
	)
	return err
}

func (s *Users) GetById(id string) (u *models.User, err error) {
	u = new(models.User)
	err = s.db.Get(
		u,
		"select * from users where id=$1",
		id,
	)
	return
}

func (s *Users) GetByEmail(email string) (u *models.User, err error) {
	u = new(models.User)
	err = s.db.Get(
		u,
		"select * from users where email=$1",
		email,
	)
	return
}

func (s *Users) GetByGithubId(id string) (u *models.User, err error) {
	u = new(models.User)
	err = s.db.Get(
		u,
		"select * from users where github_id=$1",
		id,
	)
	return
}

func (s *Users) UpdateData(id, email, avatar, username string) error {
	_, err := s.db.Exec(
		"update users set email=$1, avatar_url=$2, username=$3 where id=$4",
		email, avatar, username, id,
	)
	return err
}
