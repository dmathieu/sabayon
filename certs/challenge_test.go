package certs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPresentChallenge(t *testing.T) {
	var comChan = make(chan string)
	var errChan = make(chan error)

	challenge := &challengeProvider{
		ComChan: comChan,
		ErrChan: errChan,
	}

	go challenge.Present("example.com", "abcd", "efgh")
	assert.Equal(t, "http://example.com/.well-known/acme-challenge/efgh", <-comChan)
	comChan <- "validate"
}

func TestCleanUpChallenge(t *testing.T) {
	var comChan = make(chan string)
	var errChan = make(chan error)

	challenge := &challengeProvider{
		ComChan: comChan,
		ErrChan: errChan,
	}

	err := challenge.CleanUp("example.com", "abcd", "efgh")
	assert.Nil(t, err)
}
