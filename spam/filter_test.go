package spam

import (
	"io"
	"strings"
	"testing"

	"github.com/zen-en-tonal/mtw/session"
)

func createMail(message string) io.Reader {
	header := "From: alice<alice@mail.com>\nTo: bob<bob@mail.com>\nSubject: Subject\n\n"
	return strings.NewReader(header + message)
}

func TestRcptMismatchFilter_Spam(t *testing.T) {
	session := session.New(
		session.WithFilters(RcptMismatchFilter()),
	)
	if err := session.SetMail("alice<alice@mail.com>"); err != nil {
		t.Error(err)
	}
	if err := session.SetRcpt("spam<smap@mail.com>"); err != nil {
		t.Error(err)
	}
	if err := session.SetData(createMail("<strong>hello</strong>")); err != nil {
		t.Error(err)
	}
	if err := session.Commit(); err == nil {
		t.Error("should error")
	}
}

func TestRcptMismatchFilter_Not_Spam(t *testing.T) {
	session := session.New(
		session.WithFilters(RcptMismatchFilter()),
	)
	if err := session.SetMail("alice<alice@mail.com>"); err != nil {
		t.Error(err)
	}
	if err := session.SetRcpt("bob<bob@mail.com>"); err != nil {
		t.Error(err)
	}
	if err := session.SetData(createMail("<strong>hello</strong>")); err != nil {
		t.Error(err)
	}
	if err := session.Commit(); err != nil {
		t.Error(err)
	}
}
