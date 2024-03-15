package transac

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func NewDBConn() *sql.DB {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=example dbname=postgres sslmode=disable")
	if err != nil {
		fmt.Println("MAMA HERE")
		fmt.Println(err)
		panic(err)
	}

	return db
}
