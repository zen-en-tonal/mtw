package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zen-en-tonal/mtw/database"
	"github.com/zen-en-tonal/mtw/webhook"
)

type webhookService struct {
	create func(bp webhook.Blueprint) (*webhook.Webhook, error)
	find   func(id webhook.WebhookID) (*webhook.Webhook, error)
	all    func() (*[]webhook.Webhook, error)
}

type webhookRoute struct {
	webhookService
	Logger
}

type webhookJson struct {
	ID          string `json:"id"`
	Endpoint    string `json:"endpoint" binding:"required"`
	Auth        string `json:"auth"`
	Schema      string `json:"schema"`
	Method      string `json:"method"`
	ContentType string `json:"content_type"`
}

func (f webhookJson) into() webhook.Blueprint {
	return webhook.Blueprint{
		ID:          f.ID,
		Endpoint:    f.Endpoint,
		Auth:        f.Auth,
		Schema:      f.Schema,
		Method:      f.Method,
		ContentType: f.ContentType,
	}
}

func (r webhookRoute) register(e *gin.Engine) {
	e.POST("/webhook", r.new)
	e.GET("/webhook/:id", r.findOne)
	e.GET("/webhooks", r.findAll)
}

func (w webhookRoute) new(c *gin.Context) {
	var form webhookJson
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	webhook, err := w.create(form.into())
	if err != nil {
		w.Logger.Error("New", "error", err, "form", form)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": webhook.ID().String()})
}

func (w webhookRoute) findOne(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	webhook, err := w.find(webhook.WebhookID(id))
	if errors.Is(err, database.ErrNotFound) {
		c.JSON(http.StatusNotFound, nil)
		return
	}
	if err != nil {
		w.Logger.Error("New", "error", err, "id", id.String())
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	bp := webhook.IntoBlueprint()
	c.JSON(http.StatusOK, webhookJson{
		ID:          bp.ID,
		Endpoint:    bp.Endpoint,
		Auth:        bp.Auth,
		Schema:      bp.Schema,
		Method:      bp.Method,
		ContentType: bp.ContentType,
	})
}

func (w webhookRoute) findAll(c *gin.Context) {
	webhooks, err := w.all()
	if err != nil {
		w.Logger.Error("findAll", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	bps := make([]webhookJson, len(*webhooks))
	for i, webhook := range *webhooks {
		bp := webhook.IntoBlueprint()
		bps[i] = webhookJson{
			ID:          bp.ID,
			Endpoint:    bp.Endpoint,
			Auth:        bp.Auth,
			Schema:      bp.Schema,
			Method:      bp.Method,
			ContentType: bp.ContentType,
		}
	}
	c.JSON(http.StatusOK, bps)
}
