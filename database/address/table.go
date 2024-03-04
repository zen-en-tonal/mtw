package address

import (
	mailbox "github.com/zen-en-tonal/mtw/session"
)

type addressTable struct {
	Address string `db:"address"`
}

// into converts a webhookTable into a Webhook.
func (w addressTable) into() (*mailbox.Address, error) {
	return mailbox.ParseAddr(w.Address)
}
