package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zen-en-tonal/mtw/webhook"
)

type webhookService struct {
	createForGet  func(endpoint string, auth string) (*webhook.Webhook, error)
	createForPost func(endpoint string, schema string, contentType string, auth string) (*webhook.Webhook, error)
	find          func(id webhook.WebhookID) (*webhook.Webhook, error)
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

func (f webhookJson) do(r *webhookRoute) (*webhook.Webhook, error) {
	if f.Method == http.MethodPost {
		return r.createForPost(f.Endpoint, f.Schema, f.ContentType, f.Auth)
	} else {
		return r.createForGet(f.Endpoint, f.Auth)
	}
}

func (r webhookRoute) register(e *gin.Engine) {
	e.POST("/webhook", r.new)
	e.GET("/webhook/:id", r.find)
}

func (w webhookRoute) new(c *gin.Context) {
	var form webhookJson
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	webhook, err := form.do(&w)
	if err != nil {
		w.Logger.Error("New", "error", err, "form", form)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": webhook.ID()})
}

func (w webhookRoute) find(c *gin.Context) {
	// webhooks, err := form.do(&w)
	// if err != nil {
	// 	w.Logger.Error("New", "error", err, "form", form)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	// c.JSON(http.StatusCreated, gin.H{"id": webhook.ID()})
}
