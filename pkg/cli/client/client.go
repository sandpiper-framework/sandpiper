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
	"time"

	"github.com/google/uuid"

	"autocare.org/sandpiper/pkg/shared/model"
)

const apiVer = "/v1"

/* Use "github.com/ddliu/go-httpclient" for a sample client reference */

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

// Login to the sandpiper server
func (c *Client) Login(username, password string) error {
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
func (c *Client) SliceByName(sliceName string) (*sandpiper.Slice, error) {
	path := "/slices/name/" + sliceName
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	slice := new(sandpiper.Slice)
	_, err = c.do(req, slice)
	return slice, err
}

// ListGrains returns a list of grains for the supplied slice
func (c *Client) ListGrains(sliceID uuid.UUID) (*sandpiper.GrainsPaginated, error) {
	var results sandpiper.GrainsPaginated

	path := "/grains/slice/" + sliceID.String()
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	_, err = c.do(req, &results)
	return &results, err
}

// ListSlices returns a list of all slices
func (c *Client) ListSlices() (*sandpiper.SlicesPaginated, error) {
	var results sandpiper.SlicesPaginated

	// todo: add paging support as an argument
	path := "/slices"
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	_, err = c.do(req, &results)
	return &results, err
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
// any `body` must be valid json
func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: c.apiPrefix + path}
	u := c.baseURL.ResolveReference(rel)
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

// do executes the request
func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if c.debug {
		dump, err := httputil.DumpRequestOut(req, true)
		if err == nil {
			fmt.Printf("/n%s", dump)
		}
		fmt.Printf("req: %v\n\nresp: %v\n", req, resp)
	}
	// consider limits using json.NewDecoder(io.LimitReader(response.Body, SomeSaneConst)).Decode(v)
	err = json.NewDecoder(resp.Body).Decode(v)
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
