package database

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

const Driver = "sqlite3"

func Migrate(db *sql.DB, dbname string) error {
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		dbname,
		driver,
	)
	if err != nil {
		return err
	}
	return m.Up()
}

func Drop(db *sql.DB, dbname string) error {
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		dbname,
		driver,
	)
	if err != nil {
		return err
	}
	return m.Drop()
}
