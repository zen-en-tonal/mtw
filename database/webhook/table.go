package webhook

import (
	"github.com/google/uuid"
	"github.com/zen-en-tonal/mtw/webhook"
)

type webhookTable struct {
	ID          uuid.UUID `db:"id"`
	Endpoint    string    `db:"endpoint"`
	Auth        string    `db:"auth"`
	Schema      string    `db:"schema"`
	Method      string    `db:"method"`
	ContentType string    `db:"content_type"`
}

// into converts a webhookTable into a Webhook.
func (w webhookTable) into(defaults ...webhook.Option) (*webhook.Webhook, error) {
	bp := webhook.Blueprint{
		ID:          w.ID.String(),
		Endpoint:    w.Endpoint,
		Auth:        w.Auth,
		Schema:      w.Schema,
		Method:      w.Method,
		ContentType: w.ContentType,
	}
	return webhook.FromBlueprint(bp, defaults...)
}
