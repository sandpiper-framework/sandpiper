// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package command implements `sandpiper` commands (add, pull, list, ...)
package command

import (
	"fmt"
	"net/url"

	args "github.com/urfave/cli/v2"
)

type listParams struct {
	addr      *url.URL // our sandpiper server
	user      string
	password  string
	sliceName string
}

// List returns a list of all grains for a slice
func List(c *args.Context) error {
	p, err := getListParams(c)
	if err != nil {
		return err
	}

	// connect to the api server (saving token)
	api, err := Connect(p.addr, p.user, p.password)
	if err != nil {
		return err
	}

	if p.sliceName == "" {
		// if slice is empty, list all slices
		slices, err := api.ListSlices()
		if err != nil {
			return err
		}
		for _, slice := range slices {
			fmt.Printf("%v", slice)
		}
	} else {
		// todo: if slice is supplied, list all grains for that slice

	}

	return nil
}

func getListParams(c *args.Context) (*listParams, error) {
	// get sandpiper global params from config file and args
	g, err := GetGlobalParams(c)
	if err != nil {
		return nil, err
	}

	return &listParams{
		addr:      g.addr,
		user:      g.user,
		password:  g.password,
		sliceName: c.String("name"),
	}, nil
}
