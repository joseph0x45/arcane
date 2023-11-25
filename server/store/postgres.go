package store

import (
	"os"
  _ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
)

func GetPostgresPool() *sqlx.DB{
  pool, err := sqlx.Connect("postgres", os.Getenv("DB_URL"))
  if err!= nil{
    panic(err)
  }
  return pool
}
