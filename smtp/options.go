package smtp

import "github.com/zen-en-tonal/mtw/mailbox"

// WithMailboxOptions sets mailbox.Option into a smtp server.
func WithMailboxOptions(options ...mailbox.Option) Option {
	mailbox := mailbox.New(options...)
	return func(b *backend) {
		b.mailbox = mailbox
	}
}

// WithLogger sets Logger into a smtp server.
func WithLogger(logger Logger) Option {
	return func(b *backend) {
		b.logger = logger
	}
}
