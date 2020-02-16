// Copyright Auto Care Association. All rights reserved.
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

	"autocare.org/sandpiper/pkg/cli/payload"
	"autocare.org/sandpiper/pkg/shared/config"
	"autocare.org/sandpiper/pkg/shared/model"
)

type params struct {
	addr      *url.URL // our sandpiper api server
	user      string
	password  string
	sliceName string
	grainType string
	grainKey  string
	fileName  string
	prompt    bool
}

// Add attempts to add a new file-based grain to a slice
func Add(c *args.Context) error {
	var grain *sandpiper.Grain

	// save parameters in a `params` struct for easy access
	p, err := getParams(c)
	if err != nil {
		return err
	}

	// connect to the api server (saving token)
	api, err := Connect(p.addr, p.user, p.password)
	if err != nil {
		return err
	}

	// lookup sliceID by name
	slice, err := api.SliceByName(p.sliceName)
	if err != nil {
		return err
	}

	// load basic info from existing grain (if found) using alternate key
	grain, err = api.GrainExists(slice.ID, p.grainType, p.grainKey)
	if err != nil {
		return err
	}

	// if grain exists, remove it first (prompt for delete unless "noprompt" flag)
	if grain.ID != uuid.Nil {
		if p.prompt {
			grain.Display() // show what we're overwriting
			if !AllowOverwrite() {
				return errors.New("grain could not be added without overwrite")
			}
		}
		err := api.DeleteGrain(grain.ID)
		if err != nil {
			return err
		}
	}

	// encode supplied file for grain's payload
	data, err := payload.FromFile(p.fileName)
	if err != nil {
		return err
	}

	// create the new grain
	grain = &sandpiper.Grain{
		SliceID:  &slice.ID,
		Type:     p.grainType,
		Key:      p.grainKey,
		Source:   filepath.Base(p.fileName),
		Encoding: "gzipb64",
		Payload:  data,
	}

	// finally, add the new grain
	return api.Add(grain)
}

func getParams(c *args.Context) (*params, error) {
	// check for required file argument
	if c.NArg() != 1 {
		return nil, fmt.Errorf("missing filename argument (see 'sandpiper --help')")
	}

	// get sandpiper api server address from config file
	cfgPath := c.String("config")
	if cfgPath == "" {
		cfgPath = DefaultConfigFile
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return nil, err
	}

	addr, err := url.Parse(cfg.Server.URL)
	if err != nil {
		return nil, err
	}

	return &params{
		addr:      addr,
		user:      c.String("user"),
		password:  c.String("password"),
		sliceName: c.String("name"),
		grainType: c.String("type"),
		grainKey:  c.String("key"),
		fileName:  c.Args().Get(0),
		prompt:    !c.Bool("noprompt"), // avoid double negative
	}, nil
}
