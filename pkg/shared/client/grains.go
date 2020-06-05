// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

// Package client is used to communicate with the api server
package client

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"sandpiper/pkg/shared/model"
)

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
