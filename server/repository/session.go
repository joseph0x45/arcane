package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/joseph0x45/arcane/server/models"
)

type SessionRepo struct {
	db *sqlx.DB
}

func NewSessionRepo(db *sqlx.DB) *SessionRepo {
	return &SessionRepo{db: db}
}

func (r *SessionRepo) Insert(session *models.Session) error {
	const query = `
    insert into sessions (
      id, user_id, is_valid
    )
    values(
      :id, :user_id, :is_valid
    )
  `
	_, err := r.db.NamedExec(query, session)
	if err != nil {
		return fmt.Errorf("Error while inserting session: %w", err)
	}
	return nil
}

func (r *SessionRepo) GetByID(id string) (*models.Session, error) {
	const query = `
    select * from sessions where id=$1 and is_valid=true
  `
	session := &models.Session{}
	err := r.db.Get(session, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("Error while getting session by ID: %w", err)
	}
	return session, nil
}

func (r *SessionRepo) Invalidate(id string) error {
	const query = "update sessions set is_valid=false where id=$1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("Error while invalidating session: %w", err)
	}
	return nil
}
