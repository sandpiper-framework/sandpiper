// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package command implements the sync commands
package command

// sync show

import (
	"fmt"
	"net/url"

	"github.com/google/uuid"
	args "github.com/urfave/cli/v2"

	"autocare.org/sandpiper/pkg/shared/model"
	"autocare.org/sandpiper/pkg/sync/client"
)

type showParams struct {
	addr     *url.URL // our sandpiper server
	user     string
	password string
	slice    string // optional (empty means show slices)
	sliceID  uuid.UUID
	full     bool
	debug    bool
}

func getShowParams(c *args.Context) (*showParams, error) {
	// get sandpiper global params from config file and args
	g, err := GetGlobalParams(c)
	if err != nil {
		return nil, err
	}
	slice := c.String("slice")
	sliceID, _ := uuid.Parse(slice) // ignore error because it might be a slice-name

	return &showParams{
		addr:     g.addr,
		user:     g.user,
		password: g.password,
		full:     c.Bool("full"),
		slice:    slice,
		sliceID:  sliceID,
		debug:    g.debug,
	}, nil
}

type showCmd struct {
	*showParams
	api *client.Client
}

// newShowCmd initiates a show command
func newShowCmd(c *args.Context) (*showCmd, error) {
	p, err := getShowParams(c)
	if err != nil {
		return nil, err
	}
	// connect to the api server (saving token)
	api, err := Connect(p.addr, p.user, p.password, p.debug)
	if err != nil {
		return nil, err
	}
	return &showCmd{showParams: p, api: api}, nil
}

func (cmd *showCmd) allSlices() error {
	result, err := cmd.api.ListSlices()
	if err != nil {
		return err
	}
	for _, slice := range result.Slices {
		// show slice information
		cmd.printSlice(&slice)
	}
	return nil
}

func (cmd *showCmd) oneSlice() error {
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
	// show slice information
	cmd.printSlice(slice)
	return nil
}

func (cmd *showCmd) printSlice(slice *sandpiper.Slice) {
	if cmd.full {
		fmt.Printf(
			"%s\nName: %s (%s)\nHash: %s\nGrains: %d\n",
			slice.ID.String(), slice.Name, slice.SliceType, slice.ContentHash, slice.ContentCount,
		)
		fmt.Printf("Metadata: %v\n", slice.Metadata)
		fmt.Printf("Subscriptions: %v\n\n", slice.Companies)
	} else {
		fmt.Printf("%s (%s) \"%s\" Grains: %d\n", slice.Name, slice.ID.String(), slice.SliceType, slice.ContentCount)
	}
}

// Show displays a list of all slices that will be synced
func Show(c *args.Context) error {
	show, err := newShowCmd(c)
	if err != nil {
		return err
	}

	if show.slice == "" {
		return show.allSlices()
	} else {
		return show.oneSlice()
	}
}
