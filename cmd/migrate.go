package main

import (
	"log"

	"database/sql"

	_ "github.com/lib/pq"
	"github.com/zen-en-tonal/mtw/database"
)

func main() {
	db, err := sql.Open(
		"postgres",
		"postgres://postgres:postgres@db:5432/postgres?sslmode=disable",
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := database.Migrate(db); err != nil {
		log.Fatal(err)
	}
}
