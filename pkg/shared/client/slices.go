// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

// Package client is used to communicate with the api server
package client

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/sandpiper-framework/sandpiper/pkg/shared/model"
)

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
