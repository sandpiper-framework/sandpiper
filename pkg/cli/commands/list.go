// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package command implements `sandpiper` commands (add, pull, list, ...)
package command

import (
	sandpiper "autocare.org/sandpiper/pkg/shared/model"
	"fmt"
	"net/url"

	args "github.com/urfave/cli/v2"
)

type listParams struct {
	addr      *url.URL // our sandpiper server
	user      string
	password  string
	sliceName string
	full      bool
	debug     bool
}

// List returns a list of all grains for a slice
func List(c *args.Context) error {
	p, err := getListParams(c)
	if err != nil {
		return err
	}

	// connect to the api server (saving token)
	api, err := Connect(p.addr, p.user, p.password, p.debug)
	if err != nil {
		return err
	}

	if p.sliceName == "" {
		// if slice is empty, list all slices
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
		// todo: if slice is supplied, list all grains for that slice
		result, err := api.ListGrains(p.sliceName)
		if err != nil {
			for i, grain := range result.Grains {
				if p.full {
					printGrainFull(i, &grain)
				} else {
					printGrainBrief(&grain)
				}
			}
		}
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
		full:      c.Bool("full"),
		debug:     g.debug,
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

func printGrainFull(i int, grain *sandpiper.Grain) {
	fmt.Printf("%d: %s\n", i, grain.ID.String())
}

func printGrainBrief(grain *sandpiper.Grain) {
	fmt.Printf("%s\n", grain.ID.String())
}
