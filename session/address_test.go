package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewAddr(t *testing.T) {
	addr, err := NewAddr("alice", "localhost.lan")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "alice@localhost.lan", addr.String())
}

func Test_NewRandomAddr(t *testing.T) {
	_, err := RandomAddr("localhost.lan")
	if err != nil {
		t.Error(err)
	}
}
