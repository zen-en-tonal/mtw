package session

import "github.com/zen-en-tonal/mtw/sync"

type filters []Filter

func (fs filters) Validate(t Transaction) error {
	functions := make([]func(Transaction) error, len(fs))
	for _, f := range fs {
		functions = append(functions, f.Validate)
	}
	return sync.TryAll(t, functions...)
}

// WithFilters sets one or more filters into Session.
// Each filters execute asynchronously.
// Returns an error immediately if execution of at least one function fails.
func WithFilters(xs ...Filter) Option {
	return func(s *Session) {
		xs = append(xs, nullFilter{})
		s.Filter = filters(xs)
	}
}

type hooks struct {
	inner    []Hook
	strategy func(arg Transaction, funcs ...func(Transaction) error) error
}

func (h hooks) Send(t Transaction) error {
	functions := make([]func(Transaction) error, len(h.inner))
	for _, f := range h.inner {
		functions = append(functions, f.Send)
	}
	return h.strategy(t, functions...)
}

// WithHooksAll sets one or more hooks into Session.
// Each hooks execute asynchronously.
// Returns an error immediately if execution of at least one function fails.
func WithHooksAll(xs ...Hook) Option {
	return func(s *Session) {
		xs = append(xs, nullHook{})
		s.Hook = hooks{
			inner:    xs,
			strategy: sync.TryAll[Transaction],
		}
	}
}

// WithHooksSome sets one or more hooks into Session.
// Each hooks execute asynchronously.
func WithHooksSome(xs ...Hook) Option {
	return func(s *Session) {
		xs = append(xs, nullHook{})
		s.Hook = hooks{
			inner:    xs,
			strategy: sync.TrySome[Transaction],
		}
	}
}

// WithLogger sets a Logger into the Session.
func WithLogger(logger Logger) Option {
	return func(s *Session) {
		s.logger = logger
	}
}
