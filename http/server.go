package http

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/zen-en-tonal/mtw/database/address"
	"github.com/zen-en-tonal/mtw/database/webhook"
	"github.com/zen-en-tonal/mtw/session"
	w "github.com/zen-en-tonal/mtw/webhook"
)

type Logger interface {
	Error(msg string, args ...any)
}

func SetRoutes(r *gin.Engine, db *sql.DB, domain string, logger Logger) {
	addrRouter := addressRoute{
		addressService{
			create:       address.Create(db, domain).WithUser,
			createRandom: address.Create(db, domain).WithRandom,
			getAll:       address.Find(db).All,
			getHooks:     webhook.NewFind(db).ByAddr,
			createHook: func(addr session.Address, id w.WebhookID) error {
				return webhook.NewRegistry(db, addr).Create(id)
			},
			removeHook: func(addr session.Address, id w.WebhookID) error {
				return webhook.NewRegistry(db, addr).Remove(id)
			},
		},
		logger,
	}
	webhookRouter := webhookRoute{
		webhookService{
			create: webhook.NewCreate(db).FromBlueprint,
			find:   webhook.NewFind(db).ByID,
			all:    webhook.NewFind(db).All,
		},
		logger,
	}

	addrRouter.register(r)
	webhookRouter.register(r)
}
