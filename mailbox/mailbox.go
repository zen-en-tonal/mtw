package mailbox

import (
	"log/slog"

	"github.com/zen-en-tonal/mtw/session"
)

type Option func(*Mailbox)

type Mailbox struct {
	filters []session.Filter
	hooks   []session.Hook
	logger  session.Logger
}

// New returns the Mailbox by using Options.
func New(options ...Option) Mailbox {
	mb := Mailbox{
		filters: nil,
		hooks:   nil,
		logger:  slog.Default(),
	}
	for _, opt := range options {
		opt(&mb)
	}
	return mb
}

// NewSession returns the configured Session.
func (m Mailbox) NewSession() session.Session {
	s := session.New(
		session.WithLogger(m.logger),
		session.WithFilters(m.filters...),
		session.WithHooksSome(m.hooks...),
	)
	return s
}
