package address

import (
	"database/sql"
	"fmt"

	"github.com/zen-en-tonal/mtw/session"
	"github.com/zen-en-tonal/mtw/sync"
)

type FindHandle struct {
	addressRepository
}

// Find returns a handle to get Addresses.
func Find(db *sql.DB) FindHandle {
	return FindHandle{newRepository(db)}
}

// All returns an array of Address.
func (f FindHandle) All() (*[]session.Address, error) {
	tables, err := f.all()
	if err != nil {
		return nil, err
	}
	addrs := make([]session.Address, len(*tables))
	for i, table := range *tables {
		addr, err := table.into()
		if err != nil {
			return nil, err
		}
		addrs[i] = *addr
	}
	return &addrs, nil
}

// Exists returns the addr exists in the DB or not.
func (f FindHandle) Exists(addr session.Address) bool {
	if _, err := f.findOne(addr.String()); err != nil {
		return false
	}
	return true
}

func (f FindHandle) Validate(t session.Transaction) error {
	validate := func(selector func(t session.Transaction) string) func(t session.Transaction) error {
		return func(t session.Transaction) error {
			maybeAddr := selector(t)
			addr, err := session.ParseAddr(maybeAddr)
			if err != nil {
				return err
			}
			if !f.Exists(*addr) {
				return fmt.Errorf(
					"addr %s is not found: %w",
					addr.String(),
					session.ErrValidation,
				)
			}
			return nil
		}
	}
	return sync.TryAll(
		t,
		validate(func(t session.Transaction) string { return t.RcptAddress() }),
		validate(func(t session.Transaction) string { return t.To() }),
	)
}
