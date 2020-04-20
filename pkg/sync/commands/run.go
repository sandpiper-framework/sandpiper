// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package command implements the sync commands
package command

// sync run

import (
	"fmt"
	"net/url"

	"github.com/google/uuid"
	args "github.com/urfave/cli/v2" // conflicts with one of our package names

	"autocare.org/sandpiper/pkg/shared/model"
	"autocare.org/sandpiper/pkg/sync/client"
)

type runParams struct {
	addr     *url.URL // our sandpiper server
	user     string
	password string
	slice    string // required
	sliceID  uuid.UUID
	noupdate bool
	debug    bool
}

func getRunParams(c *args.Context) (*runParams, error) {
	// get global params from config file and args
	g, err := GetGlobalParams(c)
	if err != nil {
		return nil, err
	}
	slice := c.String("slice")
	sliceID, _ := uuid.Parse(slice) // ignore error because it might be a slice-name

	return &runParams{
		addr:     g.addr,
		user:     g.user,
		password: g.password,
		slice:    slice,
		sliceID:  sliceID,
		noupdate: c.Bool("noupdate"),
		debug:    g.debug,
	}, nil
}

type runCmd struct {
	*runParams
	api *client.Client
}

// newRunCmd initiates a run command
func newRunCmd(c *args.Context) (*runCmd, error) {
	p, err := getRunParams(c)
	if err != nil {
		return nil, err
	}
	// connect to the api server (saving token)
	api, err := Connect(p.addr, p.user, p.password, p.debug)
	if err != nil {
		return nil, err
	}
	return &runCmd{runParams: p, api: api}, nil
}

func (cmd *runCmd) allSlices() error {
	result, err := cmd.api.ListSlices()
	if err != nil {
		return err
	}
	for _, slice := range result.Slices {
		// todo: sync the slice
		fmt.Printf("%s\n", slice.ID.String())
	}
	return nil
}

func (cmd *runCmd) oneSlice() error {
	var slice *sandpiper.Slice
	var err error

	if cmd.sliceID == uuid.Nil {
		// use provided slice-name to get the slice
		slice, err = cmd.api.SliceByName(cmd.slice)
		if err != nil {
			return err
		}
	} else {
		// use provided slice-id to get the slice
		slice, err = cmd.api.SliceByID(cmd.sliceID)
		if err != nil {
			return err
		}
	}

	// todo: sync the slice
	fmt.Printf("%s\n", slice.ID.String())

	return nil
}

// Run initiates the sync process on one or all slices
func Run(c *args.Context) error {

	run, err := newRunCmd(c)
	if err != nil {
		return err
	}

	if run.slice == "" {
		return run.allSlices()
	} else {
		return run.oneSlice()
	}
}
