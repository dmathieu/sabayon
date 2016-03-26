package heroku

import "fmt"

// SetConfigVars sets the requires config vars for challenge validation
func (s *Client) SetConfigVars(appName, key, token string) error {
	var body = struct {
		Key   string `json:"ACME_KEY"`
		Token string `json:"ACME_TOKEN"`
	}{key, token}

	var res struct{}
	return s.Patch(&res, fmt.Sprintf("/apps/%s/config-vars", appName), body)
}
