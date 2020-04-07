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

	"autocare.org/sandpiper/pkg/cli/client"
	"autocare.org/sandpiper/pkg/cli/payload"
	"autocare.org/sandpiper/pkg/shared/model"
)

// grainKey is always the same for level-1 grains
const grainKey = "level-1"

type addParams struct {
	addr      *url.URL // our sandpiper server
	user      string
	password  string
	sliceName string
	fileName  string
	prompt    bool
	debug     bool
}

// Add attempts to add a new file-based grain to a slice
func Add(c *args.Context) error {

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

	// lookup sliceID by name
	slice, err := api.SliceByName(p.sliceName)
	if err != nil {
		return err
	}

	// remove the old grain first if it exists
	err = removeExistingGrain(api, p.prompt, slice.ID, grainKey)
	if err != nil {
		return err
	}

	// encode supplied file for grain's payload
	data, err := payload.FromFile(p.fileName)
	if err != nil {
		return err
	}

	// create the new grain
	grain := &sandpiper.Grain{
		SliceID:  &slice.ID,
		Key:      grainKey,
		Source:   filepath.Base(p.fileName),
		Encoding: "z64",
		Payload:  data,
	}

	// finally, add the new grain
	return api.Add(grain)
}

func getAddParams(c *args.Context) (*addParams, error) {
	// check for required file argument
	if c.NArg() != 1 {
		return nil, fmt.Errorf("missing filename argument (see 'sandpiper --help')")
	}

	// get sandpiper global params from config file and args
	g, err := GetGlobalParams(c)
	if err != nil {
		return nil, err
	}

	return &addParams{
		addr:      g.addr,
		user:      g.user,
		password:  g.password,
		sliceName: c.String("slice"),
		fileName:  c.Args().Get(0),
		prompt:    !c.Bool("noprompt"), // avoid double negative
		debug:     g.debug,
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
