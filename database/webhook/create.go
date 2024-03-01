package webhook

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/zen-en-tonal/mtw/webhook"
)

type Create struct{ webhookRepository }

// NewCreate returns a handle to create and persist a Webhook.
func NewCreate(db *sql.DB) Create {
	return Create{newRepository(db)}
}

func (c Create) persist(table webhookTable) (*webhook.Webhook, error) {
	hook, err := table.into()
	if err != nil {
		return nil, err
	}
	if err := c.upsert(table); err != nil {
		return nil, err
	}
	return hook, nil
}

// ForGet creates and persist a Webhook to send a GET request.
func (c Create) ForGet(endpoint string, auth string) (*webhook.Webhook, error) {
	table := webhookTable{
		ID:       uuid.New(),
		Endpoint: endpoint,
		Auth:     auth,
		Method:   http.MethodGet,
	}
	return c.persist(table)
}

// ForPost creates and persist a Webhook to send a POST request.
func (c Create) ForPost(endpoint string, schema string, contentType string, auth string) (*webhook.Webhook, error) {
	table := webhookTable{
		ID:          uuid.New(),
		Endpoint:    endpoint,
		Auth:        auth,
		Method:      http.MethodPost,
		Schema:      schema,
		ContentType: contentType,
	}
	return c.persist(table)
}
