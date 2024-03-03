package session

import "time"

// WithFilters sets one or more filters into Session.
// Each filters execute asynchronously.
// Returns an error immediately if execution of at least one function fails.
func WithFilters(xs ...Filter) Option {
	return func(s *Session) {
		s.Filter = Filters(xs)
	}
}

// WithHooksAll sets one or more hooks into Session.
// Each hooks execute asynchronously.
// Returns an error immediately if execution of at least one function fails.
func WithHooksAll(xs ...Hook) Option {
	return func(s *Session) {
		s.Hook = HooksAll(xs)
	}
}

// WithHooksSome sets one or more hooks into Session.
// Each hooks execute asynchronously.
func WithHooksSome(xs ...Hook) Option {
	return func(s *Session) {
		s.Hook = HooksSome(xs)
	}
}

// WithLogger sets a Logger into the Session.
func WithLogger(logger Logger) Option {
	return func(s *Session) {
		s.logger = logger
	}
}

func WithTimeout(d time.Duration) Option {
	return func(s *Session) {
		s.timeout = d
	}
}
