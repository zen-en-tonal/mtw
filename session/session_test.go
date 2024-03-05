package session

import (
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func createMail(message string) io.Reader {
	header := "From: alice<alice@mail.com>\nTo: bob<bob@mail.com>\nSubject: Subject\n\n"
	return strings.NewReader(header + message)
}

type spyHook struct {
	res Transaction
}

func (h *spyHook) Send(t Transaction) error {
	h.res = t
	return nil
}

func TestOk(t *testing.T) {
	spy := spyHook{}
	session := New(
		AppendHooks(&spy),
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
	assert.NotEmpty(t, spy.res.ID)
	assert.Equal(t, "alice@mail.com", spy.res.sender.String())
	assert.Equal(t, "bob@mail.com", spy.res.rcpt.String())
}

type errFilter struct{}

func (f errFilter) Validate(t Transaction) error {
	return errors.New("")
}

func TestValidation(t *testing.T) {
	session := New(
		AppendFilters(errFilter{}),
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
	if err := session.Commit(); err == nil {
		t.Error("should fails")
	}
}

type tooLongHook struct{}

func (t tooLongHook) Send(_ Transaction) error {
	time.Sleep(time.Second * 10)
	return nil
}

func TestTimeout(t *testing.T) {
	session := New(
		AppendHooks(tooLongHook{}),
		WithTimeout(time.Second),
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
	if err := session.Commit(); err == nil {
		t.Error("should fails")
	}
}
