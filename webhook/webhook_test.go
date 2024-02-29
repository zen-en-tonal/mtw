package webhook

import (
	"html/template"
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

func TestTemplate(t *testing.T) {
	tmp := "sender: {{.SenderAddress}}, rcpt: {{.RcptAddress}}, subject: {{.Subject}}, text: {{.Text}}"
	tmpl, err := template.New("").Parse(tmp)
	if err != nil {
		t.Error(err)
	}

	wh := New("http://example.local", WithPost(*tmpl, ContentTypeJson))
	req, err := wh.PrepareRequest(testTransaction())
	if err != nil {
		t.Error(err)
	}

	buf := make([]byte, req.ContentLength)
	req.Body.Read(buf)

	assert.Equal(t, "sender: alice@mail.com, rcpt: bob@mail.com, subject: Subject, text: hello", string(buf))
}

func TestGet(t *testing.T) {
	wh := New("http://example.local")
	_, err := wh.PrepareRequest(testTransaction())
	if err != nil {
		t.Error(err)
	}
}
