package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/joseph0x45/arcane/server/models"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Insert(user *models.User) error {
	const query = `
    insert into users (
      id, email, password
    )
    values (
      :id, :email, :password
    );
  `
	_, err := r.db.NamedExec(query, user)
	if err != nil {
		return fmt.Errorf("Error while insertig user: %w", err)
	}
	return nil
}

func (r *UserRepo) GetByEmail(email string) (*models.User, error) {
	const query = `
    select * from users where email=$1
  `
	user := models.User{}
	err := r.db.Get(&user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("Error while getting user by email: %w", err)
	}
	return &user, nil
}

func (r *UserRepo) GetByID(id string) (*models.User, error) {
	const query = `
    select * from users where id=$1
  `
	user := models.User{}
	err := r.db.Get(&user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("Error while getting user by id: %w", err)
	}
	return &user, nil
}
