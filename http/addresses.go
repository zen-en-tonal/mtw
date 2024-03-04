package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zen-en-tonal/mtw/session"
	"github.com/zen-en-tonal/mtw/webhook"
)

type addressService struct {
	create       func(user string) (*session.Address, error)
	createRandom func() (*session.Address, error)
	getAll       func() (*[]session.Address, error)
	getHooks     func(addr session.Address) (*[]webhook.Webhook, error)
	createHook   func(addr session.Address, id webhook.WebhookID) error
	removeHook   func(addr session.Address, id webhook.WebhookID) error
}

type addressRoute struct {
	addressService
	Logger
}

func (r addressRoute) register(e *gin.Engine) {
	e.GET("/addresses", r.all)
	e.POST("/address/user/random", r.newRandom)
	e.POST("/address/user/:user", r.new)
	e.GET("/address/:addr/webhooks", r.hooks)
	e.POST("/address/:addr/webhook/:whid", r.newHook)
	e.DELETE("/address/:addr/webhook/:whid", r.deleteHook)
}

func (a addressRoute) new(c *gin.Context) {
	addr, err := a.create(c.Param("user"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"address": addr.String()})
}

func (a addressRoute) newRandom(c *gin.Context) {
	addr, err := a.createRandom()
	if err != nil {
		a.Logger.Error("New", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"address": addr.String()})
}

func (a addressRoute) all(c *gin.Context) {
	addrs, err := a.getAll()
	if err != nil {
		a.Logger.Error("All", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := make([]string, len(*addrs))
	for i, a := range *addrs {
		res[i] = a.String()
	}

	c.JSON(http.StatusOK, gin.H{"addresses": res})
}

func (a addressRoute) hooks(c *gin.Context) {
	addr, err := session.ParseAddr(c.Param("addr"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hooks, err := a.getHooks(*addr)
	if err != nil {
		a.Logger.Error("hooks", "error", err, "addr", addr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ids := make([]string, len(*hooks))
	for i, hook := range *hooks {
		ids[i] = hook.ID().String()
	}

	c.JSON(http.StatusOK, gin.H{"webhooks": ids})
}

func (r addressRoute) newHook(c *gin.Context) {
	addr, err := session.ParseAddr(c.Param("addr"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hookID, err := uuid.Parse(c.Param("whid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := r.createHook(*addr, webhook.WebhookID(hookID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}

func (r addressRoute) deleteHook(c *gin.Context) {
	addr, err := session.ParseAddr(c.Param("addr"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hookID, err := uuid.Parse(c.Param("whid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := r.removeHook(*addr, webhook.WebhookID(hookID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
