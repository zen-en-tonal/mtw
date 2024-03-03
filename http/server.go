package http

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/zen-en-tonal/mtw/database/address"
	"github.com/zen-en-tonal/mtw/database/webhook"
	"github.com/zen-en-tonal/mtw/mailbox"
	w "github.com/zen-en-tonal/mtw/webhook"
)

type Logger interface {
	Error(msg string, args ...any)
}

func NewWithDB(db *sql.DB, domain string, logger Logger) *gin.Engine {
	addrRouter := addressRoute{
		addressService{
			create:       address.Create(db, domain).WithUser,
			createRandom: address.Create(db, domain).WithRandom,
			getAll:       address.Find(db).All,
			getHooks:     webhook.NewFind(db).ByAddr,
			createHook: func(addr mailbox.Address, id w.WebhookID) error {
				return webhook.NewRegistry(db, addr).Create(id)
			},
			removeHook: func(addr mailbox.Address, id w.WebhookID) error {
				return webhook.NewRegistry(db, addr).Remove(id)
			},
		},
		logger,
	}
	webhookRouter := webhookRoute{
		webhookService{
			createForGet:  webhook.NewCreate(db).ForGet,
			createForPost: webhook.NewCreate(db).ForPost,
			find:          webhook.NewFind(db).ByID,
		},
		logger,
	}

	router := gin.Default()
	addrRouter.register(router)
	webhookRouter.register(router)

	return router
}