package mailbox

import (
	"github.com/zen-en-tonal/mtw/session"
)

type FilterSet interface {
	FindFilters(addr Address) ([]session.Filter, error)
}

type filterSet struct{ FilterSet }

func (f filterSet) Validate(trans session.Transaction) error {
	addr, err := ParseAddr(trans.Envelope.GetHeader("To"))
	if err != nil {
		return err
	}
	filters, err := f.FindFilters(*addr)
	if err != nil {
		return err
	}
	return session.Filters(filters).Validate(trans)
}
