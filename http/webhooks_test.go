package http

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/zen-en-tonal/mtw/webhook"
)

func newWebhooksRoute(s webhookService) webhookRoute {
	return webhookRoute{
		webhookService: s,
		Logger:         slog.Default(),
	}
}

func Test_POST_Webhook_ForGet(t *testing.T) {
	id := uuid.MustParse("271be94b-36d1-802e-d200-c1e0b85580b2")
	router := gin.Default()
	newWebhooksRoute(webhookService{
		createForGet: func(endpoint, auth string) (*webhook.Webhook, error) {
			w := webhook.New(endpoint, webhook.WithID(id))
			return &w, nil
		},
	}).register(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"POST",
		"/webhook",
		strings.NewReader(`{"method":"GET","endpoint":"http://endpoint.com"}`),
	)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, `{"id":"271be94b-36d1-802e-d200-c1e0b85580b2"}`, w.Body.String())
}

func Test_POST_Webhook_ForPost(t *testing.T) {
	id := uuid.MustParse("271be94b-36d1-802e-d200-c1e0b85580b2")
	router := gin.Default()
	newWebhooksRoute(webhookService{
		createForPost: func(endpoint, schema, contentType, auth string) (*webhook.Webhook, error) {
			w := webhook.New(endpoint, webhook.WithID(id))
			return &w, nil
		},
	}).register(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"POST",
		"/webhook",
		strings.NewReader(`{"method":"POST","endpoint":"http://endpoint.com"}`),
	)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, `{"id":"271be94b-36d1-802e-d200-c1e0b85580b2"}`, w.Body.String())
}

func Test_GET_Webhook(t *testing.T) {
	id := uuid.MustParse("271be94b-36d1-802e-d200-c1e0b85580b2")
	router := gin.Default()
	newWebhooksRoute(webhookService{
		find: func(_ webhook.WebhookID) (*webhook.Webhook, error) {
			w := webhook.New("http://endpoint.com", webhook.WithID(id))
			return &w, nil
		},
	}).register(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"GET",
		"/webhook/271be94b-36d1-802e-d200-c1e0b85580b2",
		nil,
	)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"id":"271be94b-36d1-802e-d200-c1e0b85580b2"}`, w.Body.String())
}
