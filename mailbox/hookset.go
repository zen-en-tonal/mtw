package mailbox

import (
	"github.com/zen-en-tonal/mtw/session"
)

type HookSet interface {
	FindHooks(addr Address) ([]session.Hook, error)
}

type hookSet struct{ HookSet }

func (h hookSet) Send(trans session.Transaction) error {
	addr, err := ParseAddr(trans.Envelope.GetHeader("To"))
	if err != nil {
		return err
	}
	hooks, err := h.FindHooks(*addr)
	if err != nil {
		return err
	}
	return session.HooksSome(hooks).Send(trans)
}
