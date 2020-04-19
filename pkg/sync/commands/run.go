// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package command implements the sync commands
package command

// sync run

import (
	"net/url"

	"github.com/google/uuid"
	args "github.com/urfave/cli/v2" // conflicts with one of our package names

)

type addParams struct {
	addr     *url.URL // our sandpiper server
	user     string
	password string
	slice    string // required
	sliceID  uuid.UUID
	fileName string
	prompt   bool
	debug    bool
}

// Run performs the sync operation
func Run(c *args.Context) error {

	// save parameters in a `params` struct for easy access
	p, err := getAddParams(c)
	if err != nil {
		return err
	}

	// connect to the api server (saving token)
	api, err := Connect(p.addr, p.user, p.password, p.debug)
	if err != nil {
		return err
	}

	if p.sliceID == uuid.Nil {
		// lookup sliceID by name
		slice, err := api.SliceByName(p.slice)
		if err != nil {
			return err
		}
		p.sliceID = slice.ID
	}

	// todo: sync p.sliceID

	return nil
}

func getAddParams(c *args.Context) (*addParams, error) {

	// get sandpiper global params from config file and args
	g, err := GetGlobalParams(c)
	if err != nil {
		return nil, err
	}

	slice := c.String("slice")
	sliceID, _ := uuid.Parse(slice)

	return &addParams{
		addr:     g.addr,
		user:     g.user,
		password: g.password,
		slice:    slice,
		sliceID:  sliceID,
		prompt:   !c.Bool("noprompt"), // avoid double negative
		debug:    g.debug,
	}, nil
}
