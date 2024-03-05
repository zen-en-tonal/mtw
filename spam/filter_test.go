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

func prepareSession(options ...session.Option) session.Session {
	session := session.New(
		options...,
	)
	if err := session.SetMail("alice<alice@mail.com>"); err != nil {
		panic(err)
	}
	if err := session.SetRcpt("bob<bob@mail.com>"); err != nil {
		panic(err)
	}
	if err := session.SetData(createMail("<strong>hello</strong>")); err != nil {
		panic(err)
	}
	return session
}

func TestRcptMismatchFilter_Spam(t *testing.T) {
	session := prepareSession(
		session.AppendFilters(RcptMismatchFilter()),
	)
	if err := session.SetRcpt("tom<tom@mail.com>"); err != nil {
		panic(err)
	}
	if err := session.Commit(); err == nil {
		t.Error("should error")
	}
}

func TestRcptMismatchFilter_Not_Spam(t *testing.T) {
	session := prepareSession(
		session.AppendFilters(RcptMismatchFilter()),
	)
	if err := session.SetRcpt("bob<bob@mail.com>"); err != nil {
		panic(err)
	}
	if err := session.Commit(); err != nil {
		t.Error(err)
	}
}

func TestBlacklist_Reject(t *testing.T) {
	session := prepareSession(
		session.AppendFilters(BlackListFilter(
			`^apple@[a-z]+\.[a-z]+$`,
			`^spam@[a-z]+\.[a-z]+$`,
		)),
	)
	if err := session.SetRcpt("spam<spam@mail.com>"); err != nil {
		panic(err)
	}
	if err := session.Commit(); err == nil {
		t.Error("should error")
	}
}

func TestBlacklist_Pass(t *testing.T) {
	session := prepareSession(
		session.AppendFilters(BlackListFilter(
			`^spam@[a-z]+\.[a-z]+$`,
		)),
	)
	if err := session.SetRcpt("bob<bob@mail.com>"); err != nil {
		panic(err)
	}
	if err := session.Commit(); err != nil {
		t.Error(err)
	}
}
