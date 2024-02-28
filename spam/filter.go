package spam

import (
	"fmt"
	"net/mail"

	"github.com/zen-en-tonal/mtw/session"
)

// RcptMismatchFilter returns a filter that compares `rcpt` and `to`.
func RcptMismatchFilter() rcptMismatchFilter {
	return rcptMismatchFilter{}
}

type rcptMismatchFilter struct{}

func (r rcptMismatchFilter) Validate(e session.Transaction) error {
	rcpt := e.Rcpt.Address
	to, err := mail.ParseAddress(e.Envelope.GetHeader("To"))
	if err != nil {
		return session.ErrNilEnvelope
	}
	if rcpt != to.Address {
		return fmt.Errorf("rcpt %s and to %s is mismatched", rcpt, to.Address)
	}
	return nil
}
