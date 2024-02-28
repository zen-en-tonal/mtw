package mailbox

import "github.com/zen-en-tonal/mtw/session"

func WithFilters(filters ...session.Filter) Option {
	return func(m *Mailbox) {
		m.filters = append(m.filters, filters...)
	}
}

func WithFilterSet(set FilterSet) Option {
	return func(m *Mailbox) {
		m.filters = append(m.filters, filterSet{set})
	}
}

func WithHooks(hooks ...session.Hook) Option {
	return func(m *Mailbox) {
		m.hooks = append(m.hooks, hooks...)
	}
}

func WithHookSet(set HookSet) Option {
	return func(m *Mailbox) {
		m.hooks = append(m.hooks, hookSet{set})
	}
}

func WithLogger(logger session.Logger) Option {
	return func(m *Mailbox) {
		m.logger = logger
	}
}
