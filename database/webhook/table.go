package webhook

import (
	"html/template"
	"net/http"

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
func (w webhookTable) into() (*webhook.Webhook, error) {
	var options []webhook.Option
	options = append(options, webhook.WithID(w.ID))
	if w.Method == http.MethodPost {
		tmpl, err := template.New("").Parse(w.Schema)
		if err != nil {
			return nil, err
		}
		options = append(options, webhook.WithPost(*tmpl, w.ContentType))
	}
	if w.Auth != "" {
		options = append(options, webhook.WithAuth(w.Auth))
	}
	wh := webhook.New(w.Endpoint, options...)
	return &wh, nil
}
