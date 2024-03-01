package webhook

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/zen-en-tonal/mtw/mailbox"
)

const Driver = "postgres"

type WebhookRepository struct {
	conn *sqlx.DB
}

func NewRepository(db *sql.DB) WebhookRepository {
	return WebhookRepository{sqlx.NewDb(db, Driver)}
}

func (r WebhookRepository) upsert(table webhookTable) error {
	_, err := r.conn.Exec(`
		INSERT INTO webhooks VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id)
		DO
		UPDATE SET
			endpoint = $2
		,	auth = $3
		,	schema = $4
		,	method = $5
		,	content_type = $6
		`,
		table.ID,
		table.Endpoint,
		table.Auth,
		table.Schema,
		table.Method,
		table.ContentType,
	)
	return err
}

func (r WebhookRepository) findOne(id uuid.UUID) (*webhookTable, error) {
	table := new(webhookTable)
	if err := r.conn.Select(&table, `
		SELECT
			webhooks.*
		FROM
			webhooks
		WHERE
			id = $1
		`,
		id); err != nil {
		return nil, err
	}
	if table == nil {
		return nil, errors.New("")
	}
	return table, nil
}

func (r WebhookRepository) findByAddr(addr mailbox.Address) (*[]webhookTable, error) {
	var tables []webhookTable
	if err := r.conn.Select(&tables, `
		SELECT
			webhooks.*
		FROM
			webhooks
			JOIN
				addresses_webhooks ON webhooks.id = addresses_webhooks.webhook_id
		WHERE
			addresses_webhooks.address = $1
		`,
		addr.String()); err != nil {
		return nil, err
	}
	return &tables, nil
}

func (r *WebhookRepository) insertAddressWebhook(addr mailbox.Address, webhookID uuid.UUID) error {
	_, err := r.conn.Exec(`
		INSERT INTO addresses_webhooks (
			address
		, 	webhook_id
		)
		VALUES (
			$1
		, 	$2
		)`,
		addr.String(),
		webhookID)
	return err
}

func (r *WebhookRepository) deleteAddressWebhook(addr mailbox.Address, webhookID uuid.UUID) error {
	_, err := r.conn.Exec(`
		DELETE FROM addresses_webhooks
		WHERE
			address = $1
		AND webhook_id = $2
		`,
		addr.String(),
		webhookID)
	return err
}