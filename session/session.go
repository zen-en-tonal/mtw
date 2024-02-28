package session

import (
	"io"
	"log/slog"
	"net/mail"

	"github.com/google/uuid"
	"github.com/jhillyerd/enmime"
)

// Filter determains the mail should be recieved.
type Filter interface {
	// Validate validates a Transaction.
	// Returns an error if the transaction is invalid.
	Validate(t Transaction) error
}

// nullFilter always returns nil on Validate.
type nullFilter struct{}

func (f nullFilter) Validate(t Transaction) error {
	return nil
}

// Hook hooks
type Hook interface {
	// Send sends a Transaction.
	// Returns an error if sending fails.
	Send(t Transaction) error
}

// nullHook always returns nil on Send.
type nullHook struct{}

func (h nullHook) Send(t Transaction) error {
	return nil
}

type Logger interface {
	Error(meg string, args ...any)
}

type Session struct {
	Filter
	Hook

	logger Logger

	id       uuid.UUID
	sender   *mail.Address
	rcpt     *mail.Address
	envelope *enmime.Envelope
	data     io.Reader
}

type Option func(*Session)

func New(options ...Option) Session {
	s := Session{
		Filter: nullFilter{},
		Hook:   nullHook{},
		id:     uuid.New(),
		logger: slog.Default(),
	}
	for _, opt := range options {
		opt(&s)
	}
	return s
}

// SetMail parse a sender address and sets it into the Session.
func (s *Session) SetMail(addr string) error {
	a, err := mail.ParseAddress(addr)
	if err != nil {
		return err
	}
	s.sender = a
	return nil
}

// SetRcpt parse a recipient address and sets it into the Session.
func (s *Session) SetRcpt(addr string) error {
	a, err := mail.ParseAddress(addr)
	if err != nil {
		return err
	}
	s.rcpt = a
	return nil
}

// SetData parse body into an Envelope and sets it into the Session.
func (s *Session) SetData(r io.Reader) error {
	env, err := enmime.ReadEnvelope(r)
	if err != nil {
		return err
	}
	s.envelope = env
	s.data = r
	return nil
}

// Reset sets default values into Session.
func (s *Session) Reset() {
	s.sender = nil
	s.rcpt = nil
	s.envelope = nil
}

// Commit creates, validates, and sends a Transaction.
//
// # Errors
//   - `sender`, `rcpt`, or `envelope` is nil.
//   - Validate failed.
//   - Send failed.
func (s Session) Commit() error {
	trans, err := s.intoTransaction()
	if err != nil {
		return err
	}
	if err := s.Validate(*trans); err != nil {
		s.logger.Error(
			"validation failure",
			"reason", err,
			"id", trans.ID.String(),
			"sender", trans.Sender.String(),
			"rcpt", trans.Rcpt.String(),
			"from", trans.Envelope.GetHeader("From"),
			"to", trans.Envelope.GetHeader("To"),
			"subject", trans.Envelope.GetHeader("Subject"),
			"text", trans.Envelope.Text,
		)
		return err
	}
	return s.Send(*trans)
}

func (s Session) intoTransaction() (*Transaction, error) {
	if s.envelope == nil || s.data == nil {
		return nil, ErrNilEnvelope
	}
	if s.rcpt == nil {
		return nil, ErrNilRcpt
	}
	if s.sender == nil {
		return nil, ErrNilSender
	}
	return &Transaction{
		ID:       s.id,
		Sender:   *s.sender,
		Rcpt:     *s.rcpt,
		Envelope: *s.envelope,
		Raw:      s.data,
	}, nil
}

type Transaction struct {
	ID       uuid.UUID
	Sender   mail.Address
	Rcpt     mail.Address
	Envelope enmime.Envelope
	Raw      io.Reader
}
