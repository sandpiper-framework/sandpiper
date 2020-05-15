// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package command implements `sandpiper` commands (add, pull, list, ...)
package command

// sandpiper add

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/google/uuid"
	args "github.com/urfave/cli/v2" // conflicts with our package name

	"sandpiper/pkg/cli/payload"
	"sandpiper/pkg/shared/client"
	"sandpiper/pkg/shared/model"
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

// Add attempts to add a new file-based grain to a slice
func Add(c *args.Context) error {

	// save parameters in a `params` struct for easy access
	p, err := getAddParams(c)
	if err != nil {
		return err
	}

	// connect to our api server (saving token)
	api, err := client.Login(p.addr, p.user, p.password, p.debug)
	if err != nil {
		return err
	}

	// make sure we are on a primary-server
	if api.ServerRole() != "primary" {
		return errors.New("must be a \"primary\" server for `add` command")
	}

	// make sure we have a slice_id to work with
	if p.sliceID == uuid.Nil {
		// lookup sliceID by name
		slice, err := api.SliceByName(p.slice)
		if err != nil {
			return err
		}
		p.sliceID = slice.ID
	}

	// encode supplied file for grain's payload
	data, err := payload.FromFile(p.fileName, L1Encoding)
	if err != nil {
		return err
	}

	// lock the slice
	if err := api.LockSlice(p.sliceID); err != nil {
		return err
	}

	// remove the old grain first if it exists
	if err := removeExistingGrain(api, p.prompt, p.sliceID, sandpiper.L1GrainKey); err != nil {
		return err
	}

	// create the new grain
	grain := &sandpiper.Grain{
		SliceID:  &p.sliceID,
		Key:      sandpiper.L1GrainKey,
		Source:   filepath.Base(p.fileName),
		Encoding: L1Encoding,
		Payload:  data,
	}

	// add the new grain
	if err := api.AddGrain(grain); err != nil {
		return err
	}

	// finally, update slice content information
	if err := api.RefreshSlice(p.sliceID); err != nil {
		return err
	}

	// unlock the slice here
	if err := api.UnlockSlice(p.sliceID); err != nil {
		return err
	}

	return nil
}

func getAddParams(c *args.Context) (*addParams, error) {
	// check for required file argument
	if c.NArg() != 1 {
		return nil, fmt.Errorf("missing filename argument (see 'sandpiper -u user add --help')")
	}

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
		fileName: c.Args().Get(0),
		prompt:   !c.Bool("noprompt"), // avoid double negative
		debug:    g.debug,
	}, nil
}

func removeExistingGrain(api *client.Client, prompt bool, sliceID uuid.UUID, grainKey string) error {
	// load basic info from existing grain (if found) using alternate key
	grain, err := api.GrainExists(sliceID, grainKey)
	if err != nil {
		return err
	}

	// if grain exists, must remove it (prompt for delete unless "noprompt" flag)
	if grain.ID != uuid.Nil {
		if prompt {
			fmt.Println(grain.Display()) // show what we're overwriting
			if !AllowOverwrite() {
				return errors.New("grain could not be added without overwrite")
			}
		}
		err := api.DeleteGrain(grain.ID)
		if err != nil {
			return err
		}
	}
	return nil
}
