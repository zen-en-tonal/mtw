package smtp

import (
	"github.com/zen-en-tonal/mtw/session"
)

// WithSessionOptions sets mailbox.Option into a smtp server.
func WithSessionOptions(options ...session.Option) Option {
	return func(b *backend) {
		b.options = options
	}
}

// WithLogger sets Logger into a smtp server.
func WithLogger(logger Logger) Option {
	return func(b *backend) {
		b.logger = logger
	}
}
