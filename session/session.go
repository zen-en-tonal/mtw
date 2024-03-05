package session

import (
	"bytes"
	"io"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jhillyerd/enmime"
	"github.com/zen-en-tonal/mtw/sync"
)

type FilterProvider interface {
	GetFilters(t Transaction) (*[]Filter, error)
}

// Filter determains the mail should be recieved.
type Filter interface {
	// Validate validates a Transaction.
	// Returns an error if the transaction is invalid.
	Validate(t Transaction) error
}

// Filters is an array of Filter.
type Filters []Filter

func (f Filters) GetFilters(t Transaction) (*[]Filter, error) {
	fs := []Filter(f)
	return &fs, nil
}

func (f Filters) Validate(t Transaction, l Logger) error {
	fs := make([]func(Transaction) error, len(f))
	for i, x := range f {
		fs[i] = func(t Transaction) error {
			err := x.Validate(t)
			if err != nil {
				l.Error("")
			}
			return err
		}
	}
	return sync.TryAll(t, fs...)
}

// Hook hooks
type Hook interface {
	// Send sends a Transaction.
	// Returns an error if sending fails.
	Send(t Transaction) error
}

type HookProvider interface {
	GetHooks(t Transaction) (*[]Hook, error)
}

type Hooks []Hook

func (h Hooks) GetHooks(t Transaction) (*[]Hook, error) {
	hs := []Hook(h)
	return &hs, nil
}

func (f Hooks) Send(t Transaction, l Logger) error {
	fs := make([]func(Transaction) error, len(f))
	for i, x := range f {
		fs[i] = func(t Transaction) error {
			err := x.Send(t)
			if err != nil {
				l.Error("")
			}
			return err
		}
	}
	return sync.TrySome(t, fs...)
}

type Logger interface {
	Error(meg string, args ...any)
}

type Session struct {
	filterProviders []FilterProvider
	hookProviders   []HookProvider

	logger Logger

	id     uuid.UUID
	sender *Address
	rcpt   *Address
	data   io.Reader

	timeout time.Duration
}

type Option func(*Session)

func New(options ...Option) Session {
	s := Session{
		filterProviders: []FilterProvider{Filters{}},
		hookProviders:   []HookProvider{Hooks{}},
		id:              uuid.New(),
		logger:          slog.Default(),
		timeout:         time.Second * 10,
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
	a, err := ParseAddr(addr)
	if err != nil {
		return err
	}
	s.sender = a
	return nil
}

// SetRcpt parse a recipient address and sets it into the Session.
func (s *Session) SetRcpt(addr string) error {
	a, err := ParseAddr(addr)
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

func (s Session) filters(t Transaction) (*Filters, error) {
	filters := Filters{}
	for _, fp := range s.filterProviders {
		fs, err := fp.GetFilters(t)
		if err != nil {
			return nil, err
		}
		filters = append(filters, *fs...)
	}
	return &filters, nil
}

func (s Session) hooks(t Transaction) (*Hooks, error) {
	hooks := Hooks{}
	for _, hp := range s.hookProviders {
		fs, err := hp.GetHooks(t)
		if err != nil {
			return nil, err
		}
		hooks = append(hooks, *fs...)
	}
	return &hooks, nil
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

		filters, err := s.filters(*trans)
		if err != nil {
			s.logger.Error("")
			ec <- err
			return
		}

		hooks, err := s.hooks(*trans)
		if err != nil {
			s.logger.Error("")
			ec <- err
			return
		}

		if err := filters.Validate(*trans, s.logger); err != nil {
			ec <- err
			return
		}
		ec <- hooks.Send(*trans, s.logger)
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
	sender   Address
	rcpt     Address
	envelope enmime.Envelope
	raw      []byte
}

func NewTransaction(id uuid.UUID, sender Address, rcpt Address, body io.Reader) (*Transaction, error) {
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
	return t.sender.Name()
}

func (t Transaction) SenderAddress() string {
	return t.sender.String()
}

func (t Transaction) RcptName() string {
	return t.rcpt.Name()
}

func (t Transaction) RcptAddress() string {
	return t.rcpt.String()
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
