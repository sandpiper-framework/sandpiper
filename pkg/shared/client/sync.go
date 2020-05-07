// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package client is used to communicate with the api server
package client

import (
	"fmt"
	"github.com/google/uuid"

	"sandpiper/pkg/shared/model"
)

// ActiveServers returns a list of syncable servers
func (c *Client) ActiveServers(companyID uuid.UUID, name string) ([]sandpiper.Company, error) {
	var results []sandpiper.Company

	path := "/servers"
	if companyID != uuid.Nil {
		path = fmt.Sprintf("/servers/%s", companyID.String())
	}
	if name != "" {
		path = path + "?name=" + name
	}
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	_, err = c.do(req, results)
	return results, err
}

// AllSubs returns a list of all information we need for a sync
func (c *Client) AllSubs() ([]sandpiper.Subscription, error) {
	var results []sandpiper.Subscription

	req, err := c.newRequest("GET", "/sync/subs", nil)
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

// SubByName returns the subscription matching the supplied name
func (c *Client) SubByName(name string) (*sandpiper.Subscription, error) {
	sub := new(sandpiper.Subscription)

	path := fmt.Sprintf("/subs/name/%s", name)
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	_, err = c.do(req, sub)
	return sub, err
}

// Sync initiates a sync with a server
func (c *Client) Sync(company sandpiper.Company) error {
	path := fmt.Sprintf("/sync/%s", company.SyncAddr)
	req, err := c.newRequest("POST", path, nil)
	if err != nil {
		return err
	}
	_, err = c.do(req, nil)
	return err
}

// Process performs receives a sync request and acts on it
func (c *Client) Process() error {
	req, err := c.newRequestWS("GET", "/sync", nil)
	if err != nil {
		return err
	}
	_, err = c.do(req, nil)
	return err
}
