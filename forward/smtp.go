package forward

import (
	"net/smtp"

	"github.com/zen-en-tonal/mtw/session"
)

type Forwarder struct {
	auth       smtp.Auth
	recipients []string
	host       string
}

func NewSmtp(host string, auth smtp.Auth, recp ...string) Forwarder {
	return Forwarder{
		auth:       auth,
		recipients: recp,
		host:       host,
	}
}

func (f Forwarder) Send(t session.Transaction) error {
	return smtp.SendMail(f.host+":587", f.auth, t.From(), f.recipients, t.Raw())
}
