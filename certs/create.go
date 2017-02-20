package certs

import (
	"crypto/rand"
	"crypto/rsa"

	"github.com/xenolf/lego/acme"
)

const (
	rsaKeySize = 2048
)

// Create creates a new certificate
func (c *Cert) Create() {
	privateKey, err := rsa.GenerateKey(rand.Reader, rsaKeySize)
	if err != nil {
		c.ErrChan <- err
		return
	}

	var u = &user{
		Email: c.Email,
		key:   privateKey,
	}
	client, err := acme.NewClient(c.AcmeServer, u, acme.RSA2048)
	if err != nil {
		c.ErrChan <- err
		return
	}

	reg, err := client.Register()
	if err != nil {
		c.ErrChan <- err
		return
	}
	u.Registration = reg

	err = client.AgreeToTOS()
	if err != nil {
		c.ErrChan <- err
		return
	}

	client.ExcludeChallenges([]acme.Challenge{acme.DNS01, acme.TLSSNI01})
	challenge := &challengeProvider{
		ComChan:       c.ComChan,
		ChallengeChan: c.ChallengeChan,
		ErrChan:       c.ErrChan,
	}
	client.SetChallengeProvider(acme.HTTP01, challenge)

	certificate, errs := client.ObtainCertificate(c.Domains, true, nil, false)
	if len(errs) > 0 {
		for _, e := range errs {
			c.ErrChan <- e
		}
	}

	c.CertChan <- certificate
}
