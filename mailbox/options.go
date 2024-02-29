package mailbox

import "github.com/zen-en-tonal/mtw/session"

// WithFilters appends filters into Mailbox.
func WithFilters(filters ...session.Filter) Option {
	return func(m *Mailbox) {
		m.filters = append(m.filters, filters...)
	}
}

// WithFilterSet appends FilterSet into Mailbox.
func WithFilterSet(set FilterSet) Option {
	return func(m *Mailbox) {
		m.filters = append(m.filters, filterSet{set})
	}
}

// WithHooks appends hooks into Mailbox.
func WithHooks(hooks ...session.Hook) Option {
	return func(m *Mailbox) {
		m.hooks = append(m.hooks, hooks...)
	}
}

// WithHookSet appends HookSet into Mailbox.
func WithHookSet(set HookSet) Option {
	return func(m *Mailbox) {
		m.hooks = append(m.hooks, hookSet{set})
	}
}

// WithLogger sets a Logger into Mailbox.
func WithLogger(logger session.Logger) Option {
	return func(m *Mailbox) {
		m.logger = logger
	}
}
