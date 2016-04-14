package heroku

import "fmt"

// SetConfigVars sets the requires config vars for challenge validation
func (s *Client) SetConfigVars(appName string, index int, key, token string) error {
	var keyConfig, tokenConfig string
	if index == 0 {
		keyConfig = "ACME_KEY"
		tokenConfig = "ACME_TOKEN"
	} else {
		keyConfig = fmt.Sprintf("ACME_KEY_%d", index)
		tokenConfig = fmt.Sprintf("ACME_TOKEN_%d", index)
	}

	var body = fmt.Sprintf(
		"{\"%s\": \"%s\",\"%s\": \"%s\"}",
		keyConfig, key,
		tokenConfig, token,
	)

	var res struct{}
	return s.Patch(&res, fmt.Sprintf("/apps/%s/config-vars", appName), body)
}
