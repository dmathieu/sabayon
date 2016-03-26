package certs

import (
	"testing"

	"github.com/xenolf/lego/acme"
)

func TestUserImplementsAcmeUser(t *testing.T) {
	var _ acme.User = (*user)(nil)
}
