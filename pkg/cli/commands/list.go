// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

// Package command implements `sandpiper` commands (add, pull, list, ...)
package command

// sandpiper list command

import (
	"fmt"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/client"
	"net/url"

	"github.com/google/uuid"
	args "github.com/urfave/cli/v2"

	"github.com/sandpiper-framework/sandpiper/pkg/shared/model"
)

type listParams struct {
	addr     *url.URL // our sandpiper server
	user     string
	password string
	slice    string // optional (empty means show slices)
	sliceID  uuid.UUID
	full     bool
	debug    bool
}

// List returns a list of all grains for a slice
func List(c *args.Context) error {
	var slice *sandpiper.Slice

	p, err := getListParams(c)
	if err != nil {
		return err
	}

	// Login to the api server (saving token)
	api, err := client.Login(p.addr, p.user, p.password, p.debug)
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

func getListParams(c *args.Context) (*listParams, error) {
	// get sandpiper global params from config file and args
	g, err := GetGlobalParams(c)
	if err != nil {
		return nil, err
	}

	slice := c.String("slice")
	sliceID, _ := uuid.Parse(slice)

	return &listParams{
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
