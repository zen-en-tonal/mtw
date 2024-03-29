package session

// FilterSet represents a set of filters.
type FilterSet interface {
	// FindFilters returns an array of Filters or an error.
	// If no Filters matched the key `addr`, returns an empty array.
	FindFilters(addr Address) ([]Filter, error)
}

type filterSet struct{ FilterSet }

func AsFilter(f FilterSet) filterSet {
	return filterSet{f}
}

func (f filterSet) Validate(trans Transaction) error {
	addr, err := ParseAddr(trans.To())
	if err != nil {
		return err
	}
	filters, err := f.FindFilters(*addr)
	if err != nil {
		return err
	}
	return Filters(filters).Validate(trans)
}
