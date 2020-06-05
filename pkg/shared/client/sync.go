// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

// Package client is used to communicate with the api server
package client

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"sandpiper/pkg/shared/model"
)

// ActiveServers returns a list of syncable servers
func (c *Client) ActiveServers(companyID uuid.UUID, name string) ([]sandpiper.Company, error) {
	var servers []sandpiper.Company

	path := "/servers"
	if companyID != uuid.Nil {
		path = fmt.Sprintf("/servers/%s", companyID)
	}
	if name != "" {
		path = path + "?name=" + name
	}
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	_, err = c.do(req, &servers)
	return servers, err
}

// AllSubs returns a list of all information we need for a sync
func (c *Client) AllSubs() ([]sandpiper.Subscription, error) {
	var results []sandpiper.Subscription

	req, err := c.newRequest("GET", "/sync/subs", nil)
	if err != nil {
		return nil, err
	}
	_, err = c.do(req, &results)
	return results, err
}

// SubsByCompany returns a list of slices for supplied company_id
func (c *Client) SubsByCompany(companyID uuid.UUID) ([]sandpiper.Subscription, error) {
	var results []sandpiper.Subscription

	// todo: add paging support (looping to retrieve everything)
	path := fmt.Sprintf("/companies/%s/subs", companyID)
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	_, err = c.do(req, &results)
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

// GrainIDs returns grain-ids for a slice
func (c *Client) GrainIDs(sliceID uuid.UUID) ([]sandpiper.Grain, error) {
	var results []sandpiper.Grain

	path := fmt.Sprintf("/sync/slice/%s?brief=yes", sliceID)
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	_, err = c.do(req, &results)
	return results, err
}

// Grain returns grain (including payload) by id
func (c *Client) Grain(grainID uuid.UUID) (*sandpiper.Grain, error) {
	results := new(sandpiper.Grain)
	path := fmt.Sprintf("/grains/%s", grainID)
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	_, err = c.do(req, results)
	return results, err
}

// SliceMetaData returns an array of slice metadata records for a slice
func (c *Client) SliceMetaData(sliceID uuid.UUID) (sandpiper.MetaArray, error) {
	var results sandpiper.MetaArray
	path := fmt.Sprintf("/slices/metadata/%s", sliceID)
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	_, err = c.do(req, &results)
	return results, err
}

// LogActivity adds an activity record to the primary server
func (c *Client) LogActivity(serverID, subID uuid.UUID, msg string, duration time.Duration, e error) error {
	errMsg := ""
	if e != nil {
		errMsg = fmt.Sprintf("%v", e)
	}
	activity := sandpiper.Activity{
		CompanyID: serverID,
		SubID:     subID,
		Success:   e == nil,
		Message:   msg,
		Error:     errMsg,
		Duration:  duration,
	}
	body, err := json.Marshal(activity)
	if err != nil {
		return err
	}
	req, err := c.newRequest("POST", "/activity", body)
	if err != nil {
		return err
	}
	_, err = c.do(req, nil)
	return err
}

// Sync initiates a sync with a primary server from secondary server
func (c *Client) Sync(company sandpiper.Company) error {
	path := fmt.Sprintf("/sync/%s", company.ID)
	req, err := c.newRequest("POST", path, nil)
	if err != nil {
		return err
	}
	_, err = c.do(req, nil)
	return err
}

// Process receives a sync request and acts on it (currently unused)
func (c *Client) Process() error {
	req, err := c.newRequestWS("GET", "/sync", nil)
	if err != nil {
		return err
	}
	_, err = c.do(req, nil)
	return err
}
