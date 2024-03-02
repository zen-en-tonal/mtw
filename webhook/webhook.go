package webhook

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/zen-en-tonal/mtw/session"
)

const (
	ContentTypeJson string = "application/json"
)

type Option func(*Webhook)

func WithPost(scheme template.Template, contentType string) Option {
	return func(w *Webhook) {
		w.header.Set("Content-Type", contentType)
		w.method = "POST"
		w.schema = &scheme
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

type Logger interface {
	Error(msg string, args ...any)
}

type WebhookID uuid.UUID

// String returns an uuid as string e.g. 271be94b-36d1-802e-d200-c1e0b85580b2
func (i WebhookID) String() string {
	return uuid.UUID(i).String()
}

type Webhook struct {
	http.Client
	id       WebhookID
	endpoint string
	method   string
	header   http.Header
	schema   *template.Template
	logger   Logger
}

func New(endpoint string, options ...Option) Webhook {
	w := Webhook{endpoint: endpoint}
	WithDefault()(&w)
	for _, opt := range options {
		opt(&w)
	}
	return w
}

func (e Webhook) ID() WebhookID {
	return e.id
}

func (w Webhook) Send(t session.Transaction) error {
	req, err := w.PrepareRequest(t)
	if err != nil {
		return err
	}
	resp, err := w.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		w.logger.Error(
			"sent an http request but it responded with an error status",
			"ID", t.ID.String(),
			"Endpoint", resp.Request.URL.String(),
			"Method", resp.Request.Method,
			"StatusCode", resp.StatusCode,
			"Status", resp.Status,
		)
		return fmt.Errorf(
			"sent an http request but it responded with an error status '%s'",
			resp.Status,
		)
	}
	return nil
}

// PrepareRequest returns the `http.Request` or an error using `session.Transaction`.
func (w Webhook) PrepareRequest(t session.Transaction) (*http.Request, error) {
	var body io.Reader = nil
	if w.schema != nil {
		r, err := execTemplate(*w.schema, t)
		if err != nil {
			return nil, err
		}
		body = r
	}
	req, err := http.NewRequest(w.method, w.endpoint, body)
	if err != nil {
		return nil, err
	}
	req.Header = w.header
	return req, nil
}

func execTemplate(tmpl template.Template, t session.Transaction) (io.Reader, error) {
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, t); err != nil {
		return nil, err
	}
	return buf, nil
}
