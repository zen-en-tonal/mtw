package smtp

import (
	"errors"
	"io"
	"log/slog"

	"github.com/emersion/go-smtp"
	"github.com/zen-en-tonal/mtw/session"
)

var (
	// Empty error to hide internal message.
	Err error = errors.New("")
)

// New returns a smtp server.
func New(options ...Option) *smtp.Server {
	backend := backend{
		logger: slog.Default(),
	}
	for _, opt := range options {
		opt(&backend)
	}
	return smtp.NewServer(backend)
}

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}

type Option func(*backend)

type backend struct {
	logger  Logger
	options []session.Option
}

func (b backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	s := session.New(b.options...)
	return &smtpSession{s, b.logger}, nil
}

type smtpSession struct {
	inner  session.Session
	logger Logger
}

func (s *smtpSession) AuthPlain(username, password string) error {
	s.logger.Info("AUTH", "msg", "someone try to login", "session_id", s.inner.ID())
	return smtp.ErrAuthUnsupported
}

func (s *smtpSession) Mail(from string, opts *smtp.MailOptions) error {
	s.logger.Info("MAIL", "from", from, "session_id", s.inner.ID())
	if err := s.inner.SetMail(from); err != nil {
		s.logger.Error("MAIL", "inner", err, "from", from, "session_id", s.inner.ID())
		return Err
	}
	return nil
}

func (s *smtpSession) Rcpt(to string, opts *smtp.RcptOptions) error {
	s.logger.Info("RCPT", "to", to, "session_id", s.inner.ID())
	if err := s.inner.SetRcpt(to); err != nil {
		s.logger.Error("RCPT", "inner", err, "to", to, "session_id", s.inner.ID())
		return Err
	}
	return nil
}

func (s *smtpSession) Data(r io.Reader) error {
	s.logger.Info("DATA", "session_id", s.inner.ID())
	s.inner.SetData(r)
	if err := s.inner.Commit(); err != nil {
		s.logger.Error("DATA", "inner", err, "session_id", s.inner.ID())
		return Err
	}
	return nil
}

func (s *smtpSession) Reset() {
	s.logger.Info("RESET", "session_id", s.inner.ID())
	s.inner.Reset()
}

func (s *smtpSession) Logout() error {
	s.logger.Info("QUIT", "session_id", s.inner.ID())
	s.inner.Reset()
	return nil
}
