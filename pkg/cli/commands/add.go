// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package command

import (
	"fmt"
	"net/url"
	"os"

	args "github.com/urfave/cli/v2"

	"autocare.org/sandpiper/pkg/cli/client"
	"autocare.org/sandpiper/pkg/shared/model"
)

type params struct {
	addr      *url.URL
	user      string
	password  string
	sliceName string
	grainType string
	grainKey  string
	fileName  string
}

// Add attempts to add a new file-based grain to a slice
func Add(c *args.Context) error {

	// parse parameters
	p, err := getParams(c)
	if err != nil {
		return err
	}

	// connect to the api server
	api, err := connect(p)
	if err != nil {
		return err
	}

	// lookup sliceID by name (should we add an api to add by name to avoid this extra lookup?)
	slice, err := api.SliceByName(p.sliceName)
	if err != nil{
		return err
	}

	// todo: delete existing grain (if slice-id, grain-type, grain-key is found)
	// prompt for delete unless "noprompt" flag

	// get a reader for the file to add
	file, err := os.Open(p.fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// encode file contents for grain's payload
	payload, err := sandpiper.Encode(file)
	if err != nil {
		return err
	}

	// create the new grain
	grain := &sandpiper.Grain{
		SliceID:  &slice.ID,
		Type:     p.grainType,
		Key:      p.grainKey,
		Source:   p.fileName,
		Encoding: "gzipb64",
		Payload:  payload,
	}

	return api.Add(grain)
}

func getParams(c *args.Context) (*params, error) {
	// check for required file argument
	if c.NArg() != 1 {
		return nil, fmt.Errorf("missing filename argument (see 'sandpiper --help')")
	}

	// get sandpiper server address from command line
	addr, err := url.Parse(c.String("url"))
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
	}, nil
}

func connect(p *params) (*client.Client, error) {
	http := client.New(p.addr)
	if err := http.Login(p.user, p.password); err != nil {
		return nil, err
	}
	return http, nil
}
