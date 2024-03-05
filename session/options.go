package session

import "time"

// AppendFilters sets one or more filters into Session.
// Each filters execute asynchronously.
// Returns an error immediately if execution of at least one function fails.
func AppendFilters(xs ...Filter) Option {
	return func(s *Session) {
		s.filterProviders = append(s.filterProviders, Filters(xs))
	}
}

func AppendFilterProviders(p ...FilterProvider) Option {
	return func(s *Session) {
		s.filterProviders = append(s.filterProviders, p...)
	}
}

// WithHooksAll sets one or more hooks into Session.
// Each hooks execute asynchronously.
// Returns an error immediately if execution of at least one function fails.
func AppendHooks(xs ...Hook) Option {
	return func(s *Session) {
		s.hookProviders = append(s.hookProviders, Hooks(xs))
	}
}

func AppendHookProviders(p ...HookProvider) Option {
	return func(s *Session) {
		s.hookProviders = append(s.hookProviders, p...)
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
