package session

// HookSet represents a set of hooks.
type HookSet interface {
	// FindHooks returns an array of Hooks or an error.
	// If no Hooks matched the key `addr`, returns an empty array.
	FindHooks(addr Address) ([]Hook, error)
}

type hookSet struct{ HookSet }

func AsHook(h HookSet) hookSet {
	return hookSet{h}
}

func (h hookSet) Send(trans Transaction) error {
	addr, err := ParseAddr(trans.To())
	if err != nil {
		return err
	}
	hooks, err := h.FindHooks(*addr)
	if err != nil {
		return err
	}
	return HooksSome(hooks).Send(trans)
}
