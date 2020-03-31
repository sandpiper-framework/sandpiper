// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package client is used to communicate with the api server
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"

	"autocare.org/sandpiper/pkg/shared/model"
)

/* Use "github.com/ddliu/go-httpclient" for a sample client reference */

// Client represents the http client
type Client struct {
	BaseURL    *url.URL // basePath holds the path to prepend to the requests.
	UserAgent  string
	Auth       *sandpiper.AuthToken
	httpClient *http.Client // client used to send and receive http requests.
}

// New creates a new http client for the given sandpiper server url
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

// Login to the sandpiper server
func (c *Client) Login(username, password string) error {
	body := fmt.Sprintf("{\"username\": \"%s\", \"password\": \"%s\"}", username, password)
	req, err := c.newRequest("POST", "/login", body)
	if err != nil {
		return err
	}
	_, err = c.do(req, c.Auth)
	return err
}

// Add a grain to a slice with option to overwrite
func (c *Client) Add(grain *sandpiper.Grain) error {
	body, err := json.Marshal(grain)
	if err != nil {
		return err
	}
	req, err := c.newRequest("POST", "/grains", body)
	if err != nil {
		return err
	}
	_, err = c.do(req, nil)
	return err
}

// SliceByName returns a slice by unique key name
func (c *Client) SliceByName(name string) (*sandpiper.Slice, error) {
	path := "/slices/name/" + name
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	slice := new(sandpiper.Slice)
	_, err = c.do(req, slice)
	return slice, err
}

// ListSlices returns a list of all slices
func (c *Client) ListSlices() ([]sandpiper.Slice, error) {
	var slices []sandpiper.Slice

	// todo: add paging support as an argument
	path := "/slices"
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	_, err = c.do(req, slices)
	return slices, err
}

// GrainExists will return basic information about a grain if it exists
func (c *Client) GrainExists(sliceID uuid.UUID, grainKey string) (*sandpiper.Grain, error) {
	path := fmt.Sprintf("/grains/%s/%s", sliceID.String(), grainKey)
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	grain := new(sandpiper.Grain)
	_, err = c.do(req, grain)
	return grain, err
}

// DeleteGrain deletes a grain by primary key
func (c *Client) DeleteGrain(grainID uuid.UUID) error {
	path := fmt.Sprintf("/grains/%s", grainID.String())
	req, err := c.newRequest("DELETE", path, nil)
	if err != nil {
		return err
	}
	grain := new(sandpiper.Grain)
	_, err = c.do(req, grain)
	return err
}

/* Utility Routines */

// newRequest prepares a request for an api call
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
	if c.Auth != nil {
		req.Header.Set("Authorization", "Bearer "+c.Auth.Token)
	}
	return req, nil
}

// do executes the request
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
