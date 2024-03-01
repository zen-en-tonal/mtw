package webhook

import (
	"github.com/google/uuid"
	"github.com/zen-en-tonal/mtw/mailbox"
	"github.com/zen-en-tonal/mtw/webhook"
)

type Find struct{ WebhookRepository }

func (f Find) ByAddr(addr mailbox.Address) (*[]webhook.Webhook, error) {
	tables, err := f.findByAddr(addr)
	if err != nil {
		return nil, err
	}
	hooks := make([]webhook.Webhook, len(*tables))
	for _, table := range *tables {
		hook, err := table.into()
		if err != nil {
			return nil, err
		}
		hooks = append(hooks, *hook)
	}
	return &hooks, nil
}

func (f Find) ByID(id uuid.UUID) (*webhook.Webhook, error) {
	table, err := f.findOne(id)
	if err != nil {
		return nil, err
	}
	hook, err := table.into()
	if err == nil {
		return nil, err
	}
	return hook, nil
}
