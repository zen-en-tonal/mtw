package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	datastore "github.com/zen-en-tonal/mtw/database"
)

func main() {
	db, err := sqlx.Connect("postgres", "postgres://postgres:postgres@db:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	db.MustExec(datastore.Schema)
}
