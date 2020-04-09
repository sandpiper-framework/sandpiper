// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package command implements `sandpiper` commands (add, pull, list, ...)
package command

import (
	"fmt"
	"net/url"

	"github.com/google/uuid"
	args "github.com/urfave/cli/v2"

	"autocare.org/sandpiper/pkg/shared/model"
)

type pullParams struct {
	addr     *url.URL // our sandpiper server
	user     string
	password string
	argument string
	slice    string // optional (empty means all slices)
	sliceID  uuid.UUID
	debug    bool
}

/*
sandpiper [global-options] pull [command-options] <root-directory>
   --slice value, -s value  either a slice_id (uuid) or slice_name (case-insensitive)
*/

// Pull saves file-based grains to the file system
func Pull(c *args.Context) error {
	var slice *sandpiper.Slice

	p, err := getPullParams(c)
	if err != nil {
		return err
	}

	// connect to the api server (saving token)
	api, err := Connect(p.addr, p.user, p.password, p.debug)
	if err != nil {
		return err
	}

	if p.slice == "" {
		// no slice provided means to use all slices
		result, err := api.ListSlices()
		if err != nil {
			return err
		}
		for _, slice := range result.Slices {
			grain, err := api.GetLevel1Grain(slice.ID)
			if err != nil {
				return err
			}
			if err:=saveGrainToFile(p.argument, slice.Name, grain.Source, grain.Key, grain.Payload); err != nil {
				return err
			}
		}
	} else {
		if p.sliceID == uuid.Nil {
			// use provided slice-name to get the slice
			slice, err = api.SliceByName(p.slice)
			if err != nil {
				return err
			}
		} else {
			// use provided slice-id to get the slice
			slice, err = api.SliceByID(p.sliceID)
			if err != nil {
				return err
			}
		}
		grain, err := api.GetLevel1Grain(slice.ID)
		if err != nil {
			return err
		}
		if err:=saveGrainToFile(p.argument, slice.Name, grain.Source, grain.Key, grain.Payload); err != nil {
			return err
		}
	}
	return nil
}

func saveGrainToFile(basePath, sliceName, fileName, key string, payload sandpiper.PayloadData) error {
	if key != L1GrainKey {
		return fmt.Errorf("slice \"%s\" contains non L1 grains", sliceName)
	}
	fmt.Printf("Saving: %s/%s/%s ...", basePath, sliceName, fileName)
	return nil
}

func getPullParams(c *args.Context) (*pullParams, error) {
	// get sandpiper global params from config file and args
	g, err := GetGlobalParams(c)
	if err != nil {
		return nil, err
	}

	slice := c.String("slice")
	sliceID, _ := uuid.Parse(slice)

	return &pullParams{
		addr:     g.addr,
		user:     g.user,
		password: g.password,
		argument: c.Args().Get(0),
		slice:    slice,
		sliceID:  sliceID,
		debug:    g.debug,
	}, nil
}
