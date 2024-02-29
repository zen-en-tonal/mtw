package webhook

import (
	"html/template"
	"net/http"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/zen-en-tonal/mtw/mailbox"
	"github.com/zen-en-tonal/mtw/webhook"
)

type Webhook struct {
	ID          uuid.UUID `db:"id"`
	Endpoint    string    `db:"endpoint"`
	Auth        string    `db:"auth"`
	Schema      string    `db:"schema"`
	Method      string    `db:"method"`
	ContentType string    `db:"content_type"`
}

func (w Webhook) Into() (*webhook.Webhook, error) {
	var options []webhook.Option
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

type WebhookRepository struct {
	conn *sqlx.DB
}

func NewRepository(driver sqlx.DB) WebhookRepository {
	return WebhookRepository{&driver}
}

func (r *WebhookRepository) CreateForGet(endpoint string, auth string) (*Webhook, error) {
	table := Webhook{
		ID:       uuid.New(),
		Endpoint: endpoint,
		Auth:     auth,
		Method:   http.MethodGet,
	}
	if _, err := table.Into(); err != nil {
		return nil, err
	}
	if err := r.insert(table); err != nil {
		return nil, err
	}
	return &table, nil
}

func (r *WebhookRepository) CreateForPost(endpoint string, schema string, contentType string, auth string) (*Webhook, error) {
	table := Webhook{
		ID:          uuid.New(),
		Endpoint:    endpoint,
		Auth:        auth,
		Method:      http.MethodPost,
		Schema:      schema,
		ContentType: contentType,
	}
	if _, err := table.Into(); err != nil {
		return nil, err
	}
	if err := r.insert(table); err != nil {
		return nil, err
	}
	return &table, nil
}

func (r WebhookRepository) FindByAddr(addr mailbox.Address) (*[]Webhook, error) {
	var tables []Webhook
	if err := r.conn.Select(&tables,
		`
		SELECT
			webhooks.*
		FROM
			webhooks
			JOIN
				addresses_webhooks ON webhooks.id = addresses_webhooks.webhook_id
		WHERE
			addresses_webhooks.address = $1`,
		addr.String()); err != nil {
		return nil, err
	}
	return &tables, nil
}

func (r *WebhookRepository) RegisterAddress(addr mailbox.Address, webhookID uuid.UUID) error {
	_, err := r.conn.Exec(
		"INSERT INTO addresses_webhooks (address, webhook_id) VALUES ($1, $2)",
		addr.String(),
		webhookID)
	return err
}

func (r *WebhookRepository) UnRegisterAddress(addr mailbox.Address, webhookID uuid.UUID) error {
	_, err := r.conn.Exec(
		"DELETE FROM addresses_webhooks WHERE address = $1 AND webhook_id = $2)",
		addr.String(),
		webhookID)
	return err
}

func (r *WebhookRepository) insert(table Webhook) error {
	_, err := r.conn.Exec(
		`INSERT INTO webhooks VALUES ($1, $2, $3, $4, $5, $6)`,
		table.ID,
		table.Endpoint,
		table.Auth,
		table.Schema,
		table.Method,
		table.ContentType,
	)
	return err
}
