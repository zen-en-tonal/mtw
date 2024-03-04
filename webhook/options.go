package webhook

import (
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

var tmplFuncs = map[string]interface{}{
	"Limit": func(max int, s string) string {
		runes := []rune(s)
		if len(runes) <= max {
			return s
		}
		text := make([]rune, max)
		for i, c := range runes {
			if i >= max {
				return string(text)
			}
			text[i] = c
		}
		return string(text)
	},
}

func WithMethod(method string) Option {
	return func(w *Webhook) {
		w.method = method
	}
}

func WithSchema(schema string, contentType string) (Option, error) {
	s, err := template.New("").Funcs(tmplFuncs).Parse(schema)
	if err != nil {
		return nil, err
	}
	return func(w *Webhook) {
		w.header.Set("Content-Type", contentType)
		w.schema = s
	}, nil
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
