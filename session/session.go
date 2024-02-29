package session

import (
	"io"
	"log/slog"
	"net/mail"

	"github.com/google/uuid"
	"github.com/jhillyerd/enmime"
	"github.com/zen-en-tonal/mtw/sync"
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

// Filters is an array of Filter.
type Filters []Filter

func (f Filters) Validate(t Transaction) error {
	f = append(f, nullFilter{})
	fs := make([]func(Transaction) error, len(f))
	for _, x := range f {
		fs = append(fs, x.Validate)
	}
	return sync.TryAll(t, fs...)
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

func prepareHooks(hs []Hook) []func(Transaction) error {
	hs = append(hs, nullHook{})
	functions := make([]func(Transaction) error, len(hs))
	for _, f := range hs {
		functions = append(functions, f.Send)
	}
	return functions
}

type HooksAll []Hook

func (h HooksAll) Send(t Transaction) error {
	return sync.TryAll(t, prepareHooks(h)...)
}

type HooksSome []Hook

func (h HooksSome) Send(t Transaction) error {
	return sync.TrySome(t, prepareHooks(h)...)
}

type Logger interface {
	Error(meg string, args ...any)
}

type Session struct {
	Filter
	Hook

	logger Logger

	id     uuid.UUID
	sender *mail.Address
	rcpt   *mail.Address
	data   io.Reader
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
	s.data = r
	return nil
}

// Reset sets default values into Session.
func (s *Session) Reset() {
	s.sender = nil
	s.rcpt = nil
	s.data = nil
}

// Commit creates, validates, and sends a Transaction.
//
// # Errors
//   - `sender`, `rcpt`, or `envelope` is nil.
//   - Validate failed.
//   - Send failed.
func (s Session) Commit() error {
	trans, err := s.IntoTransaction()
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

func (s Session) IntoTransaction() (*Transaction, error) {
	if s.data == nil {
		return nil, ErrNilEnvelope
	}
	if s.rcpt == nil {
		return nil, ErrNilRcpt
	}
	if s.sender == nil {
		return nil, ErrNilSender
	}
	return NewTransaction(s.id, *s.sender, *s.rcpt, s.data)
}

type Transaction struct {
	ID       uuid.UUID
	Sender   mail.Address
	Rcpt     mail.Address
	Envelope enmime.Envelope
	Raw      io.Reader
}

func NewTransaction(id uuid.UUID, sender mail.Address, rcpt mail.Address, body io.Reader) (*Transaction, error) {
	env, err := enmime.ReadEnvelope(body)
	if err != nil {
		return nil, err
	}
	return &Transaction{
		ID:       id,
		Sender:   sender,
		Rcpt:     rcpt,
		Envelope: *env,
		Raw:      body,
	}, nil
}

func (t Transaction) SenderName() string {
	return t.Sender.Name
}

func (t Transaction) SenderAddress() string {
	return t.Sender.Address
}

func (t Transaction) RcptName() string {
	return t.Rcpt.Name
}

func (t Transaction) RcptAddress() string {
	return t.Rcpt.Address
}

func (t Transaction) HTML() string {
	return t.Envelope.HTML
}

func (t Transaction) Text() string {
	return t.Envelope.Text
}

func (t Transaction) From() string {
	return t.Envelope.GetHeader("From")
}

func (t Transaction) To() string {
	return t.Envelope.GetHeader("To")
}

func (t Transaction) Subject() string {
	return t.Envelope.GetHeader("Subject")
}
