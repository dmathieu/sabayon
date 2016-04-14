package certs

import "github.com/xenolf/lego/acme"

const (
	acmeProdServer    = "https://acme-v01.api.letsencrypt.org/directory"
	acmeStagingServer = "https://acme-staging.api.letsencrypt.org/directory"
)

// Cert allow creation and renewal of certs
type Cert struct {
	Email         string
	Domains       []string
	AcmeServer    string
	CertChan      chan acme.CertificateResource
	ComChan       chan string
	ChallengeChan chan *ChallengeParams
	ErrChan       chan error
}

// NewCert creates a new Cert struct
func NewCert(email string, domains []string) *Cert {
	return &Cert{
		Email:         email,
		Domains:       domains,
		AcmeServer:    acmeProdServer,
		CertChan:      make(chan acme.CertificateResource),
		ComChan:       make(chan string),
		ChallengeChan: make(chan *ChallengeParams),
		ErrChan:       make(chan error),
	}
}
