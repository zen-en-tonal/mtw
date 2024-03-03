package mailbox

import (
	"github.com/zen-en-tonal/mtw/session"
)

type Option func(*Mailbox)

type Mailbox struct {
	sessionOptions []session.Option
}

// New returns the Mailbox by using Options.
func New(options ...session.Option) Mailbox {
	mb := Mailbox{
		options,
	}
	return mb
}

// NewSession returns the configured Session.
func (m Mailbox) NewSession() session.Session {
	s := session.New(
		m.sessionOptions...,
	)
	return s
}
