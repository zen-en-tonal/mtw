package address

import (
	"database/sql"

	"github.com/zen-en-tonal/mtw/mailbox"
)

type CreateHandle struct {
	addressRepository
	domain string
}

// Create returns a handle to create and persist a Address.
func Create(db *sql.DB, domain string) CreateHandle {
	return CreateHandle{newRepository(db), domain}
}

// WithUser persists an address with the specified username.
func (c CreateHandle) WithUser(user string) (*mailbox.Address, error) {
	addr, err := mailbox.NewAddr(user, c.domain)
	if err != nil {
		return nil, err
	}
	return c.create(*addr)
}

// WithRandom persists an address with the randomized username by uuid.
func (c CreateHandle) WithRandom() (*mailbox.Address, error) {
	addr, err := mailbox.RandomAddr(c.domain)
	if err != nil {
		return nil, err
	}
	return c.create(*addr)
}

func (c CreateHandle) create(addr mailbox.Address) (*mailbox.Address, error) {
	table := addressTable{
		Address: addr.String(),
	}
	if err := c.insert(table); err != nil {
		return nil, err
	}
	return &addr, nil
}
