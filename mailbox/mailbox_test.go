package mailbox

import (
	"io"
	"strings"
	"testing"
)

func createMail(message string) io.Reader {
	header := "From: alice<alice@mail.com>\nTo: bob<bob@mail.com>\nSubject: Subject\n\n"
	return strings.NewReader(header + message)
}

func Test_Default_Ok(t *testing.T) {
	mb := New()
	session := mb.NewSession()
	if err := session.SetMail("alice<alice@mail.com>"); err != nil {
		t.Error(err)
	}
	if err := session.SetRcpt("bob<bob@mail.com>"); err != nil {
		t.Error(err)
	}
	if err := session.SetData(createMail("hello!")); err != nil {
		t.Error(err)
	}
	if err := session.Commit(); err != nil {
		t.Error(err)
	}
}
