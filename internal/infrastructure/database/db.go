package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewDBConn() *sqlx.DB {
	db, err := sqlx.Connect("postgres", "host=localhost port=5432 user=postgres password=example dbname=postgres sslmode=disable")
	if err != nil {
		fmt.Println("MAMA HERE")
		fmt.Println(err)
		panic(err)
	}

	return db
}
