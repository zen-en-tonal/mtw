package session

import (
	"bytes"
	"io"
	"log/slog"
	"net/mail"
	"time"

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
	for i, x := range f {
		fs[i] = x.Validate
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
	for i, f := range hs {
		functions[i] = f.Send
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

	timeout time.Duration
}

type Option func(*Session)

func New(options ...Option) Session {
	s := Session{
		Filter:  nullFilter{},
		Hook:    nullHook{},
		id:      uuid.New(),
		logger:  slog.Default(),
		timeout: time.Second * 10,
	}
	for _, opt := range options {
		opt(&s)
	}
	return s
}

// ID returns uuid.
func (s Session) ID() uuid.UUID {
	return s.id
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

	ec := make(chan error, 1)
	go func() {
		defer close(ec)
		if err := s.Validate(*trans); err != nil {
			s.logger.Error(
				"validation failure",
				"reason", err,
				"id", trans.ID.String(),
				"sender", trans.SenderAddress(),
				"rcpt", trans.RcptAddress(),
				"from", trans.From(),
				"to", trans.To(),
				"subject", trans.Subject(),
				"text", trans.Text(),
			)
			ec <- err
			return
		}
		ec <- s.Send(*trans)
	}()

	select {
	case err := <-ec:
		return err
	case <-time.After(s.timeout):
		return ErrTimeout
	}
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
	sender   mail.Address
	rcpt     mail.Address
	envelope enmime.Envelope
	raw      []byte
}

func NewTransaction(id uuid.UUID, sender mail.Address, rcpt mail.Address, body io.Reader) (*Transaction, error) {
	var buf bytes.Buffer
	tee := io.TeeReader(body, &buf)
	env, err := enmime.ReadEnvelope(tee)
	if err != nil {
		return nil, err
	}
	return &Transaction{
		ID:       id,
		sender:   sender,
		rcpt:     rcpt,
		envelope: *env,
		raw:      buf.Bytes(),
	}, nil
}

func (t Transaction) SenderName() string {
	return t.sender.Name
}

func (t Transaction) SenderAddress() string {
	return t.sender.Address
}

func (t Transaction) RcptName() string {
	return t.rcpt.Name
}

func (t Transaction) RcptAddress() string {
	return t.rcpt.Address
}

func (t Transaction) HTML() string {
	return t.envelope.HTML
}

func (t Transaction) Text() string {
	return t.envelope.Text
}

func (t Transaction) From() string {
	return t.envelope.GetHeader("From")
}

func (t Transaction) Raw() []byte {
	return t.raw
}

func (t Transaction) To() string {
	return t.envelope.GetHeader("To")
}

func (t Transaction) Subject() string {
	return t.envelope.GetHeader("Subject")
}
