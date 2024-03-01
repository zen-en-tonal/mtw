package webhook

import (
	"github.com/google/uuid"
	"github.com/zen-en-tonal/mtw/mailbox"
)

type Registory struct {
	WebhookRepository
	mailbox.Address
}

func (r Registory) Create(id uuid.UUID) error {
	return r.insertAddressWebhook(r.Address, id)
}

func (r Registory) Remove(id uuid.UUID) error {
	return r.deleteAddressWebhook(r.Address, id)
}
