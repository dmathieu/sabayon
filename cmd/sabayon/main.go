package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dmathieu/sabayon/certs"
	"github.com/dmathieu/sabayon/heroku"
)

func main() {
	var force bool
	flag.BoolVar(&force, "force", false, "Force updating the certificate even if it's not about to expire")
	flag.Parse()

	var domains = strings.Split(os.Getenv("ACME_DOMAIN"), ",")
	// trim potential whitespace
	for i, domain := range domains {
		domains[i] = strings.TrimSpace(domain)
	}
	var email = os.Getenv("ACME_EMAIL")
	var token = os.Getenv("HEROKU_TOKEN")
	var appName = os.Getenv("ACME_APP_NAME")
	wait, _ := strconv.Atoi(os.Getenv("RESTART_WAIT_TIME"))
	if wait == 0 {
		wait = 20
	}

	herokuClient := heroku.NewClient(nil, token)
	certificates, err := herokuClient.GetSSLCertificates(appName)
	if err != nil {
		log.Fatal(err)
	}

	if len(certificates) > 1 {
		log.Fatalf("Found %d certificate. Can only update one. Nothing done.", len(certificates))
	}

	if len(certificates) != 0 && !force {
		certExpiration, err := time.Parse(time.RFC3339, certificates[0].SslCert.ExpiresAt)
		if err != nil {
			log.Fatal(err)
		}
		now := time.Now()
		renew := certExpiration.AddDate(0, -1, 0)

		if now.Before(renew) {
			log.Printf("cert.ignore_update expires_at=\"%s\" renew_at=\"%s\"", certExpiration, renew)
			return
		}
	}

	log.Printf("cert.create email='%s' domains='%s'", email, domains)

	ce := certs.NewCert(email, domains)
	go ce.Create()

	for {
		select {
		case r := <-ce.ErrChan:
			log.Printf("%s", r)
			return
		case r := <-ce.ChallengeChan:
			log.Printf("cert.validate")
			var index int
			for i, v := range domains {
				if v == r.Domain {
					index = i + 1
				}
			}

			err := herokuClient.SetConfigVars(appName, index, r.KeyAuth, r.Token)
			if err != nil {
				log.Fatal(err)
			}

			// Wait for a few seconds so the app can restart
			time.Sleep(time.Duration(wait) * time.Second)

			ce.ComChan <- "validate"
		case r := <-ce.CertChan:
			log.Printf("cert.created")

			if len(certificates) == 0 {
				err = herokuClient.SetSSLCertificate(appName, r.Certificate, r.PrivateKey)
				if err != nil {
					log.Fatal(err)
				}

				log.Printf("cert.added")
			} else {
				err = herokuClient.UpdateSSLCertificate(appName, certificates[0].Name, r.Certificate, r.PrivateKey)
				if err != nil {
					log.Fatal(err)
				}

				log.Printf("cert.updated")
			}

			return
		}
	}
}
