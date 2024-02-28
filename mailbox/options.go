package mailbox

import "github.com/zen-en-tonal/mtw/session"

func WithFilterSet(set FilterSet) Option {
	return func(m *Mailbox) {
		m.filters = append(m.filters, filterSet{set})
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
