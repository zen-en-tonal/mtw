package webhook

import (
	"database/sql"

	"github.com/zen-en-tonal/mtw/session"
	"github.com/zen-en-tonal/mtw/webhook"
)

type Find struct {
	webhookRepository
	options []webhook.Option
}

// NewFind returns a handle to get Webhooks.
func NewFind(db *sql.DB, defaults ...webhook.Option) Find {
	return Find{newRepository(db), defaults}
}

// ByAddr returns Webhooks by Address.
func (f Find) ByAddr(addr session.Address) (*[]webhook.Webhook, error) {
	tables, err := f.findByAddr(addr)
	if err != nil {
		return nil, err
	}
	hooks := make([]webhook.Webhook, len(*tables))
	for i, table := range *tables {
		hook, err := table.into(f.options...)
		if err != nil {
			return nil, err
		}
		hooks[i] = *hook
	}
	return &hooks, nil
}

// ByID returns a Webhook by WebhookID.
//
// # Errors
//   - If no Webhook found.
func (f Find) ByID(id webhook.WebhookID) (*webhook.Webhook, error) {
	table, err := f.findOne(id)
	if err != nil {
		return nil, err
	}
	hook, err := table.into(f.options...)
	if err != nil {
		return nil, err
	}
	return hook, nil
}

// All returns an array of Webhook.
func (f Find) All() (*[]webhook.Webhook, error) {
	tables, err := f.findAll()
	if err != nil {
		return nil, err
	}
	hooks := make([]webhook.Webhook, len(*tables))
	for i, table := range *tables {
		hook, err := table.into(f.options...)
		if err != nil {
			return nil, err
		}
		hooks[i] = *hook
	}
	return &hooks, nil
}

func (f Find) FindHooks(addr session.Address) ([]session.Hook, error) {
	webhooks, err := f.ByAddr(addr)
	if err != nil {
		return nil, err
	}
	hooks := make([]session.Hook, len(*webhooks))
	for i, webhook := range *webhooks {
		hooks[i] = webhook
	}
	return hooks, nil
}
