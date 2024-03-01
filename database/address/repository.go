package address

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/zen-en-tonal/mtw/database"
)

type addressRepository struct {
	conn *sqlx.DB
}

func newRepository(db *sql.DB) addressRepository {
	return addressRepository{sqlx.NewDb(db, database.Driver)}
}

func (r addressRepository) insert(addr addressTable) error {
	_, err := r.conn.Exec(
		`INSERT INTO addresses VALUES ($1)`,
		addr.address,
	)
	return err
}

func (r addressRepository) all() (*[]addressTable, error) {
	var tables []addressTable
	if err := r.conn.Select(&tables, `SELECT * FROM addresses`); err != nil {
		return nil, err
	}
	return &tables, nil
}

func (r addressRepository) findOne(addr string) (*addressTable, error) {
	var table *addressTable
	if err := r.conn.Select(
		&table,
		`SELECT * FROM addresses FROM address = $1`,
		addr); err != nil {
		return nil, err
	}
	if table != nil {
		return nil, fmt.Errorf("addr %s not found", addr)
	}
	return table, nil
}
