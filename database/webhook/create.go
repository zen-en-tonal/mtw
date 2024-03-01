package webhook

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/zen-en-tonal/mtw/webhook"
)

type Create struct{ webhookRepository }

func NewCreate(db *sql.DB) Create {
	return Create{newRepository(db)}
}

func (c Create) create(table webhookTable) (*webhook.Webhook, error) {
	hook, err := table.into()
	if err != nil {
		return nil, err
	}
	if err := c.upsert(table); err != nil {
		return nil, err
	}
	return hook, nil
}

func (c Create) ForGet(endpoint string, auth string) (*webhook.Webhook, error) {
	table := webhookTable{
		ID:       uuid.New(),
		Endpoint: endpoint,
		Auth:     auth,
		Method:   http.MethodGet,
	}
	return c.create(table)
}

func (c Create) ForPost(endpoint string, schema string, contentType string, auth string) (*webhook.Webhook, error) {
	table := webhookTable{
		ID:          uuid.New(),
		Endpoint:    endpoint,
		Auth:        auth,
		Method:      http.MethodPost,
		Schema:      schema,
		ContentType: contentType,
	}
	return c.create(table)
}
