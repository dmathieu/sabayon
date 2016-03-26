package heroku

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
)

const (
	apiURL    = "https://api.heroku.com"
	userAgent = "sabayon"
)

// Client holds the token and methods to make calls to the Heroku API
type Client struct {
	client *http.Client
	Token  string
}

// NewClient creates a new client
func NewClient(c *http.Client, token string) *Client {
	if c == nil {
		c = http.DefaultClient
	}
	return &Client{
		client: c,
		Token:  token,
	}
}

// NewRequest generates an HTTP request, but does not perform the request.
func (s *Client) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	var (
		ctype string
		rbody io.Reader
	)

	switch t := body.(type) {
	case nil:
	case string:
		rbody = bytes.NewBufferString(t)
	case io.Reader:
		rbody = t
	default:
		v := reflect.ValueOf(body)
		if !v.IsValid() {
			break
		}
		if v.Type().Kind() == reflect.Ptr {
			v = reflect.Indirect(v)
			if !v.IsValid() {
				break
			}
		}
		j, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		rbody = bytes.NewReader(j)
		ctype = "application/json"
	}

	req, err := http.NewRequest(method, apiURL+path, rbody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.heroku+json; version=3.sni_ssl_cert")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.Token))
	if ctype != "" {
		req.Header.Set("Content-Type", ctype)
	}
	return req, nil
}

// Do sends a request and decodes the response into v.
func (s *Client) Do(v interface{}, method, path string, body interface{}, lr *ListRange) error {
	req, err := s.NewRequest(method, path, body)
	if err != nil {
		return err
	}
	if lr != nil {
		lr.SetHeader(req)
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		contents, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Status: %d - Error: %s", resp.StatusCode, contents)
	}

	switch t := v.(type) {
	case nil:
	case io.Writer:
		_, err = io.Copy(t, resp.Body)
	default:
		err = json.NewDecoder(resp.Body).Decode(v)
	}
	return err
}

// Get sends a GET request and decodes the response into v.
func (s *Client) Get(v interface{}, path string, lr *ListRange) error {
	return s.Do(v, "GET", path, nil, lr)
}

// Patch sends a Path request and decodes the response into v.
func (s *Client) Patch(v interface{}, path string, body interface{}) error {
	return s.Do(v, "PATCH", path, body, nil)
}

// Post sends a POST request and decodes the response into v.
func (s *Client) Post(v interface{}, path string, body interface{}) error {
	return s.Do(v, "POST", path, body, nil)
}

// Put sends a PUT request and decodes the response into v.
func (s *Client) Put(v interface{}, path string, body interface{}) error {
	return s.Do(v, "PUT", path, body, nil)
}

// Delete sends a DELETE request.
func (s *Client) Delete(path string) error {
	return s.Do(nil, "DELETE", path, nil, nil)
}
