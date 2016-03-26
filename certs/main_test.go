package certs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCert(t *testing.T) {
	cert := NewCert("test@example.com", "example.com")
	assert.Equal(t, "test@example.com", cert.Email)
	assert.Equal(t, "example.com", cert.Domain)
	assert.NotNil(t, cert.CertChan)
	assert.NotNil(t, cert.ComChan)
	assert.NotNil(t, cert.ErrChan)
}
