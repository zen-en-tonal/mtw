package webhook

import (
	"io"
	"net/mail"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/zen-en-tonal/mtw/session"
)

func createMail(message string) io.Reader {
	header := "From: alice<alice@mail.com>\nTo: bob<bob@mail.com>\nSubject: Subject\n\n"
	return strings.NewReader(header + message)
}

func testTransaction() session.Transaction {
	m := createMail("hello")
	t, err := session.NewTransaction(
		uuid.New(),
		mail.Address{Address: "alice@mail.com"},
		mail.Address{Address: "bob@mail.com"},
		m,
	)
	if err != nil {
		panic(err)
	}
	return *t
}

func Test_parseOption_Post(t *testing.T) {
	webhook, err := webhookTable{
		Endpoint:    "http://example.com",
		Method:      "POST",
		Auth:        "secret",
		Schema:      "{}",
		ContentType: "application/json",
	}.into()
	if err != nil {
		t.Error(err)
	}
	req, err := webhook.PrepareRequest(testTransaction())
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, req.Method, "POST")
	assert.Equal(t, req.Header.Get("Authorization"), "secret")
	assert.Equal(t, req.Header.Get("Content-Type"), "application/json")
}
