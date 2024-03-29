package webhook

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/zen-en-tonal/mtw/database"
	"github.com/zen-en-tonal/mtw/session"
	"github.com/zen-en-tonal/mtw/webhook"
)

type webhookRepository struct {
	conn *sqlx.DB
}

func newRepository(db *sql.DB) webhookRepository {
	return webhookRepository{sqlx.NewDb(db, database.Driver)}
}

func (r webhookRepository) upsert(table webhookTable) error {
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

func (r webhookRepository) findOne(id webhook.WebhookID) (*webhookTable, error) {
	var tables []webhookTable
	if err := r.conn.Select(&tables, `
		SELECT
			webhooks.*
		FROM
			webhooks
		WHERE
			id = $1
		`,
		id.String()); err != nil {
		return nil, err
	}
	if len(tables) == 0 {
		return nil, database.ErrNotFound
	}
	table := tables[0]
	return &table, nil
}

func (r webhookRepository) findAll() (*[]webhookTable, error) {
	var tables []webhookTable
	if err := r.conn.Select(&tables, `SELECT webhooks.* FROM webhooks`); err != nil {
		return nil, err
	}
	return &tables, nil
}

func (r webhookRepository) findByAddr(addr session.Address) (*[]webhookTable, error) {
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

func (r *webhookRepository) insertAddressWebhook(addr session.Address, webhookID webhook.WebhookID) error {
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
		webhookID.String())
	return err
}

func (r *webhookRepository) deleteAddressWebhook(addr session.Address, webhookID webhook.WebhookID) error {
	_, err := r.conn.Exec(`
		DELETE FROM addresses_webhooks
		WHERE
			address = $1
		AND webhook_id = $2
		`,
		addr.String(),
		webhookID.String())
	return err
}
