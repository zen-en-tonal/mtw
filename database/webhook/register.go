package webhook

import (
	"database/sql"

	"github.com/zen-en-tonal/mtw/mailbox"
	"github.com/zen-en-tonal/mtw/webhook"
)

type Registory struct {
	webhookRepository
	mailbox.Address
}

func NewRegistory(db *sql.DB, addr mailbox.Address) Registory {
	return Registory{newRepository(db), addr}
}

func (r Registory) Create(id webhook.WebhookID) error {
	return r.insertAddressWebhook(r.Address, id)
}

func (r Registory) Remove(id webhook.WebhookID) error {
	return r.deleteAddressWebhook(r.Address, id)
}
