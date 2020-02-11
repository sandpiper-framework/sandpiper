// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package command

import (
	"autocare.org/sandpiper/pkg/shared/mock"
	sandpiper "autocare.org/sandpiper/pkg/shared/model"
	"fmt"
	"net/url"
	"os"

	args "github.com/urfave/cli/v2"

	"autocare.org/sandpiper/pkg/cli/client"
)

// Add attempts to add a new file-based grain to a slice
func Add(c *args.Context) error {
	// check for required file argument
	if c.NArg() != 1 {
		fmt.Errorf("missing filename argument (see 'sandpiper --help')")
	}

	// get sandpiper server address from command line
	addr, err := url.Parse(c.String("url"))
	if err != nil {
		return err
	}

	// login to api (saving the token in the client)
	http := client.New(addr)
	if err := http.Login(c.String("user"), c.String("password")); err != nil {
  	return err
	}

	// todo: delete existing grain (if slice-id, grain-type, grain-key is found)

	// todo: lookup sliceID by c.String("slice")
	sliceID := mock.TestUUID(1)

	// add new grain (encode payload first)
	filename := c.Args().Get(0)
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	payload, err := sandpiper.Encode(file);
	if err != nil {
		return err
	}

	grain := &sandpiper.Grain {
		SliceID: &sliceID,
		Type: c.String("type"),
		Key: c.String("key"),
		Encoding: "gzipb64",
		Payload: payload,
	}

	return http.Add(grain)
}
