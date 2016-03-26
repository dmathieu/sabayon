package certs

import (
	"crypto"

	"github.com/xenolf/lego/acme"
)

type user struct {
	Email        string
	Registration *acme.RegistrationResource
	key          crypto.PrivateKey
}

func (u *user) GetEmail() string {
	return u.Email
}

func (u *user) GetRegistration() *acme.RegistrationResource {
	return u.Registration
}

func (u *user) GetPrivateKey() crypto.PrivateKey {
	return u.key
}
