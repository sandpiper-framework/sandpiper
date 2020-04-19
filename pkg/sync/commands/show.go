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
)

type cmdParams struct {
	addr     *url.URL // our sandpiper server
	user     string
	password string
	slice    string // optional (empty means show slices)
	sliceID  uuid.UUID
	full     bool
	debug    bool
}

// Show returns a list of all slices that will be synced
func Show(c *args.Context) error {
	var slice *sandpiper.Slice

	p, err := getCmdParams(c)
	if err != nil {
		return err
	}

	// connect to the api server (saving token)
	api, err := Connect(p.addr, p.user, p.password, p.debug)
	if err != nil {
		return err
	}

	if p.slice == "" {
		// no slice provided means to list all slices
		result, err := api.ListSlices()
		if err != nil {
			return err
		}
		for i, slice := range result.Slices {
			if p.full {
				printSliceFull(i, slice)
			} else {
				printSliceBrief(slice)
			}
		}
	} else {
		if p.sliceID == uuid.Nil {
			// use provided slice-name to get the slice-id
			slice, err = api.SliceByName(p.slice)
			if err != nil {
				return err
			}
			p.sliceID = slice.ID
		}
		// return a list of paginated grains for the slice-id
		// todo: add pagination logic
		result, err := api.ListGrains(p.sliceID, p.full)
		if err != nil {
			return err
		}
		for _, grain := range result.Grains {
			if p.full {
				printGrainFull(&grain)
			} else {
				printGrainBrief(&grain)
			}
		}
	}
	return nil
}

func getCmdParams(c *args.Context) (*cmdParams, error) {
	// get sandpiper global params from config file and args
	g, err := GetGlobalParams(c)
	if err != nil {
		return nil, err
	}

	slice := c.String("slice")
	sliceID, _ := uuid.Parse(slice)

	return &cmdParams{
		addr:     g.addr,
		user:     g.user,
		password: g.password,
		full:     c.Bool("full"),
		slice:    slice,
		sliceID:  sliceID,
		debug:    g.debug,
	}, nil
}

func printSliceFull(i int, slice sandpiper.Slice) {
	fmt.Printf(
		"%d: %s\nName: %s (%s)\nHash: %s\nGrains: %d\n",
		i+1, slice.ID.String(), slice.Name, slice.SliceType, slice.ContentHash, slice.ContentCount,
	)
	fmt.Printf("Metadata: %v\n", slice.Metadata)
	fmt.Printf("Subscriptions: %v\n\n", slice.Companies)
}

func printSliceBrief(slice sandpiper.Slice) {
	fmt.Printf("%s (%s) \"%s\" Grains: %d\n", slice.Name, slice.ID.String(), slice.SliceType, slice.ContentCount)
}

func printGrainFull(grain *sandpiper.Grain) {
	fmt.Println(grain.DisplayFull())
}

func printGrainBrief(grain *sandpiper.Grain) {
	fmt.Println(grain.Display())
}
