package webhook

import (
	"database/sql"

	"github.com/zen-en-tonal/mtw/mailbox"
	"github.com/zen-en-tonal/mtw/webhook"
)

type Registry struct {
	webhookRepository
	mailbox.Address
}

// NewRegistry returns a handle to register a Webhook to the Address.
func NewRegistry(db *sql.DB, addr mailbox.Address) Registry {
	return Registry{newRepository(db), addr}
}

// Create registers the Webhook to the Address in the context.
func (r Registry) Create(id webhook.WebhookID) error {
	return r.insertAddressWebhook(r.Address, id)
}

// Remove deletes the Webhook on the Address in the context.
func (r Registry) Remove(id webhook.WebhookID) error {
	return r.deleteAddressWebhook(r.Address, id)
}
