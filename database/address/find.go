package address

import (
	"database/sql"

	"github.com/zen-en-tonal/mtw/mailbox"
)

type FindHandle struct {
	addressRepository
}

// Find returns a handle to get Addresses.
func Find(db *sql.DB) FindHandle {
	return FindHandle{newRepository(db)}
}

// All returns an array of Address.
func (f FindHandle) All() (*[]mailbox.Address, error) {
	tables, err := f.all()
	if err != nil {
		return nil, err
	}
	addrs := make([]mailbox.Address, len(*tables))
	for i, table := range *tables {
		addr, err := table.into()
		if err != nil {
			return nil, err
		}
		addrs[i] = *addr
	}
	return &addrs, nil
}
