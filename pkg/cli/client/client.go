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

//region ** Slice Routines **

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

// SliceByID returns a slice by primary key
func (c *Client) SliceByID(sliceID uuid.UUID) (*sandpiper.Slice, error) {
	path := "/slices/" + sliceID.String()
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	slice := new(sandpiper.Slice)
	_, err = c.do(req, slice)
	return slice, err
}

// ListSlices returns a list of all slices
func (c *Client) ListSlices() (*sandpiper.SlicesPaginated, error) {
	var results sandpiper.SlicesPaginated

	// todo: add paging support as an argument
	req, err := c.newRequest("GET", "/slices", nil)
	if err != nil {
		return nil, err
	}
	_, err = c.do(req, &results)
	return &results, err
}

// RefreshSlice updates the content information about a slice
func (c *Client) RefreshSlice(sliceID uuid.UUID) error {
	path := fmt.Sprintf("/slices/refresh/%s", sliceID.String())
	req, err := c.newRequest("POST", path, nil)
	if err != nil {
		return err
	}
	_, err = c.do(req, nil)
	return err
}

// LockSlice suspends sync operations from starting
func (c *Client) LockSlice(sliceID uuid.UUID) error {
	path := fmt.Sprintf("/slices/lock/%s", sliceID.String())
	req, err := c.newRequest("PUT", path, nil)
	if err != nil {
		return err
	}
	_, err = c.do(req, nil)
	return err
}

// UnlockSlice suspends sync operations from starting
func (c *Client) UnlockSlice(sliceID uuid.UUID) error {
	path := fmt.Sprintf("/slices/unlock/%s", sliceID.String())
	req, err := c.newRequest("PUT", path, nil)
	if err != nil {
		return err
	}
	_, err = c.do(req, nil)
	return err
}

//endregion

//region ** Grain Routines **

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

// GetLevel1Grain returns a grain by sliceID if Level1
func (c *Client) GetLevel1Grain(sliceID uuid.UUID) (*sandpiper.Grain, error) {
	path := fmt.Sprintf("/grains/%s/%s?payload=yes", sliceID.String(), sandpiper.L1GrainKey)
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	grain := new(sandpiper.Grain)
	resp, err := c.do(req, grain)
	if resp != nil && resp.StatusCode == 404 {
		// "not found" is not an error
		return grain, nil
	}
	return grain, err
}

// ListGrains returns a list of grains for the supplied slice
func (c *Client) ListGrains(sliceID uuid.UUID, fullFlag bool) (*sandpiper.GrainsPaginated, error) {
	var results sandpiper.GrainsPaginated
	var payloadParam string

	if fullFlag {
		payloadParam = "?payload=yes"
	}
	path := "/grains/slice/" + sliceID.String() + payloadParam
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	_, err = c.do(req, &results)
	return &results, err
}

// AddGrain adds a grain to a slice (without overwrite)
func (c *Client) AddGrain(grain *sandpiper.Grain) error {
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

// DeleteGrain deletes a grain by primary key
func (c *Client) DeleteGrain(grainID uuid.UUID) error {
	path := fmt.Sprintf("/grains/%s", grainID.String())
	req, err := c.newRequest("DELETE", path, nil)
	if err != nil {
		return err
	}
	_, err = c.do(req, nil)
	return err
}

//endregion

//region ** Sync Routines **

// AllSubs returns a list of all slices
func (c *Client) AllSubs() ([]sandpiper.Subscription, error) {
	var results []sandpiper.Subscription

	// todo: add paging support (looping to retrieve everything)
	req, err := c.newRequest("GET", "/subs", nil)
	if err != nil {
		return nil, err
	}
	_, err = c.do(req, results)
	return results, err
}

// SubsByCompany returns a list of slices for supplied company_id
func (c *Client) SubsByCompany(companyID uuid.UUID) ([]sandpiper.Subscription, error) {
	var results []sandpiper.Subscription

	// todo: add paging support (looping to retrieve everything)
	path := fmt.Sprintf("/companies/%s/subs", companyID.String())
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	_, err = c.do(req, results)
	return results, err
}

// SubsByName returns a list of slices for supplied subscription name
func (c *Client) SubsByName(name string) ([]sandpiper.Subscription, error) {
	var results []sandpiper.Subscription

	// todo: add paging support (looping to retrieve everything)
	path := fmt.Sprintf("/subs/name/%s", name)
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	_, err = c.do(req, results)
	return results, err
}

//endregion

//region ** Utility Routines **

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

//endregion
