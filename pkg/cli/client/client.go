// Copyright Auto Care Association. All rights reserved.
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

	//"github.com/ddliu/go-httpclient" use this for reference

	"github.com/google/uuid"

	"autocare.org/sandpiper/pkg/shared/model"
)

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

func (c *Client) ListUsers() ([]sandpiper.User, error) {
	req, err := c.newRequest("GET", "/users", nil)
	if err != nil {
		return nil, err
	}
	var users []sandpiper.User
	_, err = c.do(req, &users)
	return users, err
}

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

func (c *Client) GrainExists(sliceID uuid.UUID, grainType, grainKey string) (*sandpiper.Grain, error) {
	path := fmt.Sprintf("/grains/%s/%s/%s", sliceID.String(), grainType, grainKey)
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	grain := new(sandpiper.Grain)
	_, err = c.do(req, grain)
	return grain, err
}

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
