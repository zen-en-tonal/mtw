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
		`INSERT INTO addresses (address) VALUES ($1)`,
		addr.Address,
	)
	return err
}

func (r addressRepository) all() (*[]addressTable, error) {
	var tables []addressTable
	if err := r.conn.Select(&tables, `SELECT address FROM addresses`); err != nil {
		return nil, err
	}
	return &tables, nil
}

func (r addressRepository) findOne(addr string) (*addressTable, error) {
	var tables []addressTable
	if err := r.conn.Select(
		&tables,
		`SELECT address FROM addresses WHERE address = $1`,
		addr); err != nil {
		return nil, err
	}
	if len(tables) == 0 {
		return nil, fmt.Errorf("addr %s not found", addr)
	}
	return &tables[0], nil
}
