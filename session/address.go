package session

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

var mailAddressRegexp = regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)

// Address represents an email address.
type Address struct {
	user   string // e.g. Alice <`alice`@mail.com>.
	domain string // e.g. Alice <alice@`mail.com`>.
	name   string // e.g. `Alice` <alice@mail.com>. name can be empty.
}

// ParseAddr creates an Address from one line string.
//
// # Errors
//   - If address is invalid format for mail address.
func ParseAddr(address string) (*Address, error) {
	addr, err := mail.ParseAddress(address)
	if err != nil {
		return nil, err
	}
	// This line maybe not necessary
	// because mail.ParseAddress might check the address.
	if !mailAddressRegexp.Match([]byte(addr.Address)) {
		return nil, fmt.Errorf("address %s is invalid", address)
	}
	return &Address{
		user:   strings.Split(addr.Address, "@")[0],
		domain: strings.Split(addr.Address, "@")[1],
		name:   addr.Name,
	}, nil
}

// MustParseAddr creates an Address from one line string.
//
// # Panics
//   - If addr is invalid format for mail address.
func MustParseAddr(addr string) Address {
	a, err := ParseAddr(addr)
	if err != nil {
		panic(err)
	}
	return *a
}

// NewAddr creates an Address from user and domain.
//
// # Errors
//   - If provides invalid format for mail address.
func NewAddr(user string, domain string) (*Address, error) {
	return ParseAddr(user + "@" + domain)
}

// RandomAddr creates an Address with user that generated by uuid v4.
//
// # Errors
//   - If provides invalid format for mail address.
func RandomAddr(domain string) (*Address, error) {
	return NewAddr(uuid.NewString(), domain)
}

// User returns a section of user in address.
func (a Address) User() string {
	return a.user
}

// Domain returns a section of domain in address.
func (a Address) Domain() string {
	return a.domain
}

// Name returns a section of name in address.
func (a Address) Name() string {
	return a.name
}

// String returns the address as string. e.g. alice@mail.com
func (a Address) String() string {
	return a.user + "@" + a.domain
}
