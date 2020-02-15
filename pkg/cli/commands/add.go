// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package command implements `sandpiper` commands (add, pull, list, ...)
package command

// sandpiper add

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/google/uuid"
	args "github.com/urfave/cli/v2" // conflicts with our package name

	"autocare.org/sandpiper/pkg/cli/client"
	"autocare.org/sandpiper/pkg/shared/config"
	"autocare.org/sandpiper/pkg/shared/model"
)

type params struct {
	addr      *url.URL	// our sandpiper api server
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
	api, err := connect(p.addr, p.user, p.password)
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
			if !allowOverwrite() {
				return nil
			}
		}
		err := api.DeleteGrain(grain.ID)
		if err != nil {
			return err
		}
	}

	// encode supplied file for grain's payload
	payload, err := payloadFromFile(p.fileName)
	if err != nil {
		return err
	}

	// create the new grain
	grain = &sandpiper.Grain{
		SliceID:  &slice.ID,
		Type:     p.grainType,
		Key:      p.grainKey,
		Source:   p.fileName,
		Encoding: "gzipb64",
		Payload:  payload,
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

func connect(addr *url.URL, user, password string) (*client.Client, error) {
	http := client.New(addr)
	if err := http.Login(user, password); err != nil {
		return nil, err
	}
	return http, nil
}

func allowOverwrite() bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Overwrite (y/n)? ")
	ans, _ := reader.ReadString('\n')
	return strings.ToLower(ans) == "y"
}

func payloadFromFile(fileName string) (sandpiper.PayloadData, error) {
	// get a reader for the file to add
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// encode file contents for grain's payload
	payload, err := sandpiper.Encode(file)
	if err != nil {
		return nil, err
	}
	return payload, nil
}
