package storetests

import (
	"server/models"
	"server/store"
	"testing"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func TestUsersStore(t *testing.T) {
	godotenv.Load("../../.env")
	db := store.GetPostgresPool()
	store := store.NewUsersStore(db)
	user_id := uuid.NewString()
	userData := &models.User{
		Id:        user_id,
    GithubId: "random_id",
		Email:     "random@mail.com",
		Username:  "random",
		AvatarURL: "lorem picsum",
	}
	t.Run("insert new user", func(t *testing.T) {
		err := store.Insert(userData)
		if err != nil {
			t.Errorf("Wanted nil got %q", err)
		}
	})

	t.Run("get user by id", func(t *testing.T) {
		dbUser, err := store.GetById(user_id)
		if err != nil {
			t.Errorf("Wanted nil got %q", err)
		}
		if !cmp.Equal(dbUser, userData) {
			t.Errorf("Wanted %#v got %#v", userData, dbUser)
		}
	})

	t.Run("get user with email", func(t *testing.T) {
		dbUser, err := store.GetByEmail("random@mail.com")
		if err != nil {
			t.Errorf("Wanted nil got %q", err)
		}
		if !cmp.Equal(dbUser, userData) {
			t.Errorf("Wanted %#v got %#v", userData, dbUser)
		}
	})
}
