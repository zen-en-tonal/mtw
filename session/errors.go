package session

import (
	"errors"
)

var (
	ErrNilEnvelope error = errors.New("nil envelope")
	ErrNilRcpt     error = errors.New("nil rcpt")
	ErrNilSender   error = errors.New("nil sender")

	ErrValidation error = errors.New("validation failure")
	ErrTimeout    error = errors.New("timeout")
)
