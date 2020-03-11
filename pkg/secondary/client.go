// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package secondary is the sandpiper client process which polls the sandpiper server
// for changes and executes the sync.
package secondary

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"autocare.org/sandpiper/pkg/shared/model"
)

type Client struct {
	BaseURL    *url.URL // basePath holds the path to prepend to the requests.
	UserAgent  string
	httpClient *http.Client // client used to send and receive http requests.
}

func New(baseURL *url.URL) *Client {
	netClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	c := &Client{
		BaseURL:    baseURL,
		UserAgent:  "Sandpiper",
		httpClient: netClient,
	}
	return c
}

func (c *Client) ListUsers() ([]sandpiper.User, error) {
	req, err := c.newRequest("GET", "/users", nil)
	if err != nil {
		return nil, err
	}
	var users []sandpiper.User
	_, err = c.do(req, &users)
	return users, err
}

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(rel)
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// consider limits using json.NewDecoder(io.LimitReader(response.Body, SomeSaneConst)).Decode(v)
	err = json.NewDecoder(resp.Body).Decode(v)
	return resp, err
}
