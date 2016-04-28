package certs

import "log"

type challengeProvider struct {
	ComChan       chan string
	ChallengeChan chan *ChallengeParams
	ErrChan       chan error
}

// ChallengeParams holds the params we need to match for the challenge
type ChallengeParams struct {
	Domain  string
	KeyAuth string
	Token   string
}

func (c *challengeProvider) Present(domain, token, keyAuth string) error {
	c.ChallengeChan <- &ChallengeParams{
		Domain:  domain,
		KeyAuth: keyAuth,
		Token:   token,
	}

	for {
		select {
		case r := <-c.ComChan:
			if r == "validate" {
				log.Printf("cert.validated")
				return nil
			}
		}
	}
}

func (c *challengeProvider) CleanUp(domain, token, keyAuth string) error {
	return nil
}
