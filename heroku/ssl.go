package heroku

import "fmt"

// Certificate is the information of a single certificate
type Certificate struct {
	Name string `json:"name"`
}

// SetSSLCertificate adds a certificate to an app
func (s *Client) SetSSLCertificate(appName string, chain, key []byte) error {
	var body = struct {
		Chain      string `json:"certificate_chain"`
		PrivateKey string `json:"private_key"`
	}{string(chain), string(key)}

	var res struct{}
	return s.Post(&res, fmt.Sprintf("/apps/%s/sni-endpoints", appName), body)
}

// UpdateSSLCertificate updates an existing certificate
func (s *Client) UpdateSSLCertificate(appName, certName string, chain, key []byte) error {
	var body = struct {
		Chain      string `json:"certificate_chain"`
		PrivateKey string `json:"private_key"`
	}{string(chain), string(key)}

	var res struct{}
	return s.Patch(&res, fmt.Sprintf("/apps/%s/sni-endpoints/%s", appName, certName), body)
}

// GetSSLCertificates returns the certificates for an app
func (s *Client) GetSSLCertificates(appName string) ([]Certificate, error) {
	var res []Certificate
	return res, s.Get(&res, fmt.Sprintf("/apps/%s/sni-endpoints", appName), &ListRange{})
}

// RemoveSSLCertificates removes all certificates added to an app
func (s *Client) RemoveSSLCertificates(appName string) error {
	certificates, err := s.GetSSLCertificates(appName)
	if err != nil {
		return err
	}

	for _, c := range certificates {
		err = s.Delete(fmt.Sprintf("/apps/%s/sni-endpoints/%s", appName, c.Name))
		if err != nil {
			return err
		}
	}

	return nil
}
