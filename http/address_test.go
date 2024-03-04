package http

import (
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/zen-en-tonal/mtw/session"
	"github.com/zen-en-tonal/mtw/webhook"
)

func newAddrRoute(addr addressService) addressRoute {
	return addressRoute{
		addressService: addr,
		Logger:         slog.Default(),
	}
}

func TestGetAddresses(t *testing.T) {
	router := gin.Default()
	newAddrRoute(addressService{
		getAll: func() (*[]session.Address, error) {
			return &[]session.Address{
				session.MustParseAddr("alice@mail.com"),
			}, nil
		},
	}).register(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/addresses", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"addresses":["alice@mail.com"]}`, w.Body.String())
}

func TestGetAddresses_Error(t *testing.T) {
	router := gin.Default()
	newAddrRoute(addressService{
		getAll: func() (*[]session.Address, error) {
			return nil, errors.New("err")
		},
	}).register(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/addresses", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, `{"error":"err"}`, w.Body.String())
}

func TestNewAddress(t *testing.T) {
	router := gin.Default()
	newAddrRoute(addressService{
		create: func(user string) (*session.Address, error) {
			return session.NewAddr(user, "mail.com")
		},
	}).register(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"POST",
		"/address/user/alice",
		nil,
	)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, `{"address":"alice@mail.com"}`, w.Body.String())
}

func TestNewRandom(t *testing.T) {
	router := gin.Default()
	newAddrRoute(addressService{
		createRandom: func() (*session.Address, error) {
			return session.NewAddr("alice", "mail.com")
		},
	}).register(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"POST",
		"/address/user/random",
		nil,
	)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, `{"address":"alice@mail.com"}`, w.Body.String())
}

func TestHooks(t *testing.T) {
	router := gin.Default()
	wh := webhook.New("http://endpoint.com", webhook.WithID(uuid.MustParse(
		"271be94b-36d1-802e-d200-c1e0b85580b2",
	)))
	newAddrRoute(addressService{
		getHooks: func(addr session.Address) (*[]webhook.Webhook, error) {
			return &[]webhook.Webhook{wh}, nil
		},
	}).register(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"GET",
		"/address/alice@mail.com/webhooks",
		nil,
	)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"webhooks":["271be94b-36d1-802e-d200-c1e0b85580b2"]}`, w.Body.String())
}

func TestHooks_BadRequest(t *testing.T) {
	router := gin.Default()
	newAddrRoute(addressService{}).register(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"GET",
		"/address/invalid/webhooks",
		nil,
	)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func Test_POST_Hook(t *testing.T) {
	router := gin.Default()
	newAddrRoute(addressService{
		createHook: func(addr session.Address, id webhook.WebhookID) error {
			return nil
		},
	}).register(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"POST",
		"/address/alice@mail.com/webhook/271be94b-36d1-802e-d200-c1e0b85580b2",
		nil,
	)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func Test_POST_Hook_BadAddress(t *testing.T) {
	router := gin.Default()
	newAddrRoute(addressService{
		createHook: func(addr session.Address, id webhook.WebhookID) error {
			return nil
		},
	}).register(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"POST",
		"/address/alicemail.com/webhook/271be94b-36d1-802e-d200-c1e0b85580b2",
		nil,
	)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func Test_POST_Hook_BadHookID(t *testing.T) {
	router := gin.Default()
	newAddrRoute(addressService{
		createHook: func(addr session.Address, id webhook.WebhookID) error {
			return nil
		},
	}).register(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"POST",
		"/address/alice@mail.com/webhook/271be94b-36d1-802e-d200-c1e0b85",
		nil,
	)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func Test_DELETE_Hook(t *testing.T) {
	router := gin.Default()
	newAddrRoute(addressService{
		removeHook: func(addr session.Address, id webhook.WebhookID) error {
			return nil
		},
	}).register(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"DELETE",
		"/address/alice@mail.com/webhook/271be94b-36d1-802e-d200-c1e0b85580b2",
		nil,
	)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func Test_DELETE_Hook_BadAddress(t *testing.T) {
	router := gin.Default()
	newAddrRoute(addressService{
		removeHook: func(addr session.Address, id webhook.WebhookID) error {
			return nil
		},
	}).register(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"DELETE",
		"/address/alicemail.com/webhook/271be94b-36d1-802e-d200-c1e0b85580b2",
		nil,
	)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func Test_DELETE_Hook_BadHookID(t *testing.T) {
	router := gin.Default()
	newAddrRoute(addressService{
		removeHook: func(addr session.Address, id webhook.WebhookID) error {
			return nil
		},
	}).register(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"DELETE",
		"/address/alice@mail.com/webhook/271be94b-36d1-802e-d200",
		nil,
	)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
