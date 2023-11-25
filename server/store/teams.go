package store

import (
	"server/models"

	"github.com/jmoiron/sqlx"
)

type Teams struct {
  db *sqlx.DB
}

func NewTeamsStore(db *sqlx.DB) *Teams{
  return &Teams{
    db: db,
  }
}

func (s *Teams) Insert(t *models.Team) error {
  _, err := s.db.NamedExec(
    "insert into teams(id, name, owner, plan) values(:id, :name, :owner, :plan)",
    t,
  )
  return err
}

func (s *Teams) GetById(id string) (t *models.Team, err error){
  t = new(models.Team)
  err = s.db.Get(
    t,
    "select * from teams where id=$1",
    id,
  )
  return
}
