// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package client is used to communicate with the api server
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"sandpiper/pkg/shared/model"
)

const apiVer = "/v1"

// Client represents the http client
type Client struct {
	baseURL    *url.URL // basePath holds the path to prepend to the requests.
	apiPrefix  string   // prepended to endpoint after successful /login
	userAgent  string
	auth       *sandpiper.AuthToken
	httpClient *http.Client // client used to send and receive http requests.
	debug      bool
}

// New creates a new http client for the given sandpiper server url
func New(baseURL *url.URL, debugFlag bool) *Client {
	var timeout time.Duration = 10

	if debugFlag {
		timeout = 3600
	}

	netClient := &http.Client{
		Timeout: timeout * time.Second,
	}

	c := &Client{
		baseURL:    baseURL,
		userAgent:  "Sandpiper",
		auth:       &sandpiper.AuthToken{},
		httpClient: netClient,
		debug:      debugFlag,
	}
	return c
}

// Login to the sandpiper api server (saving token in the client struct)
func Login(addr *url.URL, user, password string, debug bool) (*Client, error) {
	c := New(addr, debug)
	if err := c.login(user, password); err != nil {
		return nil, err
	}
	return c, nil
}

// login to the sandpiper server
func (c *Client) login(username, password string) error {
	body := fmt.Sprintf("{\"username\": \"%s\", \"password\": \"%s\"}", username, password)
	req, err := c.newRequest("POST", "/login", body)
	if err != nil {
		return err
	}
	resp, err := c.do(req, c.auth)
	if err == nil && resp.StatusCode != 200 {
		return fmt.Errorf("login failed (%d)", resp.StatusCode)
	}

	// add api version to all subsequent api calls
	c.apiPrefix = apiVer

	return err
}

// newRequest prepares a request for an api call
// `body` (if not nil) must be valid json
func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	u, err := c.baseURL.Parse(c.apiPrefix + path)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, u.String(), toReader(body))
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	if c.auth.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.auth.Token)
	}
	return req, nil
}

// newRequestWS prepares a websocket request for an api call
// `body` (if not nil) must be valid json
func (c *Client) newRequestWS(method, path string, body interface{}) (*http.Request, error) {
	u, err := c.baseURL.Parse(c.apiPrefix + path)
	if err != nil {
		return nil, err
	}
	u.Scheme = strings.Replace(u.Scheme, "http", "ws", 1)

	req, err := http.NewRequest(method, u.String(), toReader(body))
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Connection", "upgrade")
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Sec-WebSocket-Version", "13")
	req.Header.Set("Sec-WebSocket-Key", "??what to put here??")
	if c.auth.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.auth.Token)
	}
	return req, nil
}

// do executes the request
func (c *Client) do(req *http.Request, v interface{}) (*Response, error) {
	r, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	resp := &Response{r} // wrap it in our struct for new methods
	defer resp.Body.Close()

	if c.debug {
		dump, err := httputil.DumpRequestOut(req, true)
		if err == nil {
			fmt.Printf("/n%s", dump)
		}
		fmt.Printf("req: %v\n\nresp: %v\n", req, resp)
	}

	if resp.StatusCode != 200 {
		msg, _ := resp.ToString()
		return resp, fmt.Errorf("%s: %s", resp.Status, msg)
	}

	if v != nil {
		// convert the json response to the provided structure pointer
		// consider limits using json.NewDecoder(io.LimitReader(response.Body, SomeSaneConst)).Decode(v)
		err = json.NewDecoder(resp.Body).Decode(v)
	}
	return resp, err
}

func toReader(v interface{}) *bytes.Reader {
	switch t := v.(type) {
	case []byte:
		return bytes.NewReader(t)
	case string:
		return bytes.NewReader([]byte(t))
	case *bytes.Reader:
		return t
	case nil:
		return bytes.NewReader(nil)
	default:
		panic("Invalid value")
	}
}