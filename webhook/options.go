package webhook

import (
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func WithPost(scheme template.Template, contentType string) Option {
	return func(w *Webhook) {
		w.header.Set("Content-Type", contentType)
		w.method = "POST"
		w.schema = &scheme
	}
}

func WithMethod(method string) Option {
	return func(w *Webhook) {
		w.method = method
	}
}

func WithSchema(schema template.Template, contentType string) Option {
	return func(w *Webhook) {
		w.header.Set("Content-Type", contentType)
		w.schema = &schema
	}
}

func WithAuth(token string) Option {
	return func(w *Webhook) {
		w.header.Set("Authorization", token)
	}
}

func WithTimeout(d time.Duration) Option {
	return func(w *Webhook) {
		w.Timeout = d
	}
}

func WithID(id uuid.UUID) Option {
	return func(w *Webhook) {
		w.id = WebhookID(id)
	}
}

func WithDefault() Option {
	return func(w *Webhook) {
		w.id = WebhookID(uuid.New())
		w.header = http.Header{}
		w.method = "GET"
		w.Timeout = time.Second * 10
		w.schema = nil
		w.logger = slog.Default()
	}
}
