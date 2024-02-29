package mailbox

import (
	"github.com/zen-en-tonal/mtw/session"
)

// FilterSet represents a set of filters.
type FilterSet interface {
	// FindFilters returns an array of Filters or an error.
	// If no Filters matched the key `addr`, returns an empty array.
	FindFilters(addr Address) ([]session.Filter, error)
}

type filterSet struct{ FilterSet }

func (f filterSet) Validate(trans session.Transaction) error {
	addr, err := ParseAddr(trans.To())
	if err != nil {
		return err
	}
	filters, err := f.FindFilters(*addr)
	if err != nil {
		return err
	}
	return session.Filters(filters).Validate(trans)
}
