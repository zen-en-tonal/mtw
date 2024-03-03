package mailbox

import (
	"github.com/zen-en-tonal/mtw/session"
)

// HookSet represents a set of hooks.
type HookSet interface {
	// FindHooks returns an array of Hooks or an error.
	// If no Hooks matched the key `addr`, returns an empty array.
	FindHooks(addr Address) ([]session.Hook, error)
}

type hookSet struct{ HookSet }

func AsHook(h HookSet) hookSet {
	return hookSet{h}
}

func (h hookSet) Send(trans session.Transaction) error {
	addr, err := ParseAddr(trans.To())
	if err != nil {
		return err
	}
	hooks, err := h.FindHooks(*addr)
	if err != nil {
		return err
	}
	return session.HooksSome(hooks).Send(trans)
}
