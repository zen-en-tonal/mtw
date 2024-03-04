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

func TestGet(t *testing.T) {
	wh := New("http://example.local")
	_, err := wh.PrepareRequest(testTransaction())
	if err != nil {
		t.Error(err)
	}
}

func Test_Blueprint(t *testing.T) {
	bp := Blueprint{
		Endpoint:    "http://example.local",
		Method:      "POST",
		Schema:      `{"msg":"{{.Text}}"}`,
		ContentType: "application/json",
	}
	wh, err := FromBlueprint(bp)
	if err != nil {
		t.Error(err)
	}
	req, err := wh.PrepareRequest(testTransaction())
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "http://example.local", req.URL.String())
	assert.Equal(t, "POST", req.Method)
	assert.Equal(t, "application/json", req.Header.Get("content-type"))
	buf := make([]byte, req.ContentLength)
	req.Body.Read(buf)
	assert.Equal(t, `{"msg":"hello"}`, string(buf))
}

func Test_IntoBlueprint(t *testing.T) {
	bp := Blueprint{
		ID:          "ece24b02-c98f-46b2-993f-a0860cd116cd",
		Endpoint:    "http://example.local",
		Method:      "POST",
		Schema:      `{"msg":"{{.Text}}"}`,
		ContentType: "application/json",
		Auth:        "sercret",
	}
	wh, err := FromBlueprint(bp)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, bp, wh.IntoBlueprint())
}

func Test_TemplateFunction_Limit(t *testing.T) {
	bp := Blueprint{
		Endpoint:    "http://example.local",
		Method:      "POST",
		Schema:      `{"msg":"{{Limit 1 .Text}}"}`,
		ContentType: "application/json",
	}
	wh, err := FromBlueprint(bp)
	if err != nil {
		t.Error(err)
	}
	req, err := wh.PrepareRequest(testTransaction())
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "http://example.local", req.URL.String())
	assert.Equal(t, "POST", req.Method)
	assert.Equal(t, "application/json", req.Header.Get("content-type"))
	buf := make([]byte, req.ContentLength)
	req.Body.Read(buf)
	assert.Equal(t, `{"msg":"h"}`, string(buf))
}

func Test_TemplateFunction_Limit_NoEffects(t *testing.T) {
	bp := Blueprint{
		Endpoint:    "http://example.local",
		Method:      "POST",
		Schema:      `{"msg":"{{Limit 1000 .Text}}"}`,
		ContentType: "application/json",
	}
	wh, err := FromBlueprint(bp)
	if err != nil {
		t.Error(err)
	}
	req, err := wh.PrepareRequest(testTransaction())
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "http://example.local", req.URL.String())
	assert.Equal(t, "POST", req.Method)
	assert.Equal(t, "application/json", req.Header.Get("content-type"))
	buf := make([]byte, req.ContentLength)
	req.Body.Read(buf)
	assert.Equal(t, `{"msg":"hello"}`, string(buf))
}
