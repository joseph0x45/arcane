package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/joseph0x45/arcane/models"
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
      id, github_id, username, avatar_url,
      joined_at
    )
    values (
      :id, :github_id, :username, :avatar_url,
      :joined_at
    );
  `
	_, err := r.db.NamedExec(query, user)
	if err != nil {
		return fmt.Errorf("Error while insertig user: %w", err)
	}
	return nil
}

func (r *UserRepo) UpdateUserData(username, avatarURL, userID string) error {
	const query = `
    update users set username=$1, avatar_url=$2
    where id=$3
  `
	_, err := r.db.Exec(query, username, avatarURL, userID)
	if err != nil {
		return fmt.Errorf("Error while updating user data: %w", err)
	}
	return nil
}

func (r *UserRepo) GetByGithubID(ghID string) (*models.User, error) {
	const query = `
    select * from users where github_id=$1
  `
	user := models.User{}
	err := r.db.Get(&user, query, ghID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("Error while getting user by github_id: %w", err)
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
