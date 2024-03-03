package address

import (
	"github.com/zen-en-tonal/mtw/mailbox"
)

type addressTable struct {
	Address string `db:"address"`
}

// into converts a webhookTable into a Webhook.
func (w addressTable) into() (*mailbox.Address, error) {
	return mailbox.ParseAddr(w.Address)
}
