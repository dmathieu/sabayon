package commands

import (
	"log"
	"os"
	"time"

	"github.com/codegangsta/cli"
	"github.com/dmathieu/sabayon/certs"
	"github.com/dmathieu/sabayon/heroku"
)

// SetupCmd configures a new certificate for an app
var SetupCmd = cli.Command{
	Name:  "setup",
	Usage: "Setup a cert for the configured app",

	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "day,d",
			Usage: "Only perform the task on the specified day",
		},
	},

	Action: func(c *cli.Context) {
		requiredDay := c.String("day")
		t := time.Now()
		currentDay := t.Weekday().String()

		if requiredDay != "" && requiredDay != currentDay {
			log.Printf("cert.day_ignore required=%s current=%s", requiredDay, currentDay)
			return
		}

		var domain = os.Getenv("ACME_DOMAIN")
		var email = os.Getenv("ACME_EMAIL")
		var token = os.Getenv("HEROKU_TOKEN")
		var appName = os.Getenv("ACME_APP_NAME")

		log.Printf("cert.create email='%s' domain='%s'", email, domain)

		ce := certs.NewCert(email, domain)
		go ce.Create()

		herokuClient := heroku.NewClient(nil, token)

		for {
			select {
			case r := <-ce.ErrChan:
				log.Printf("%s", r)
				return
			case r := <-ce.ChallengeChan:
				log.Printf("cert.validate")

				err := herokuClient.SetConfigVars(appName, r.KeyAuth, r.Token)
				if err != nil {
					log.Fatal(err)
				}

				// Wait for a few seconds so the app can restart
				time.Sleep(5 * time.Second)

				ce.ComChan <- "validate"
			case r := <-ce.ComChan:
				log.Printf("cert.com msg=%s", r)
			case r := <-ce.CertChan:
				log.Printf("cert.created")

				certs, err := herokuClient.GetSSLCertificates(appName)
				if err != nil {
					log.Fatal(err)
				}

				if len(certs) > 1 {
					log.Fatalf("Found %d certificate. Can only update one. Nothing done.", len(certs))
				}

				if len(certs) == 0 {
					err = herokuClient.SetSSLCertificate(appName, r.Certificate, r.PrivateKey)
					if err != nil {
						log.Fatal(err)
					}

					log.Printf("cert.added")
				} else {
					err = herokuClient.UpdateSSLCertificate(appName, certs[0].Name, r.Certificate, r.PrivateKey)
					if err != nil {
						log.Fatal(err)
					}

					log.Printf("cert.updated")
				}

				return
			}
		}

	},
}
