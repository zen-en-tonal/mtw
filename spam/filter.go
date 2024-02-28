package spam

import (
	"fmt"
	"net/mail"
	"regexp"

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

type blackList []string

// BlackListFilter sets filters defined as regexp pattern.
// If at least one pattern matches, the filter returns an error.
func BlackListFilter(patterns ...string) blackList {
	return blackList(patterns)
}

func (patterns blackList) Validate(e session.Transaction) error {
	rcpt := e.Rcpt.Address
	to, err := mail.ParseAddress(e.Envelope.GetHeader("To"))
	if err != nil {
		return session.ErrNilEnvelope
	}
	for _, pattern := range patterns {
		r, err := regexp.Compile(pattern)
		if err != nil {
			return err
		}
		if r.Match([]byte(rcpt)) {
			return fmt.Errorf("addr %s contains blacklist", rcpt)
		}
		if r.Match([]byte(to.Address)) {
			return fmt.Errorf("addr %s contains blacklist", to.Address)
		}
	}
	return nil
}
