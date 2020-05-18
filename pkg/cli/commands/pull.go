// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package command implements `sandpiper` commands (add, pull, list, ...)
package command

import (
	"fmt"
	"net/url"
	"os"

	"github.com/google/uuid"
	args "github.com/urfave/cli/v2"

	"sandpiper/pkg/shared/client"
	"sandpiper/pkg/shared/model"
)

/*
sandpiper [global-options] pull [command-options] <root-directory>
   --slice value, -s value  either a slice_id (uuid) or slice_name (case-insensitive)
*/

const folderPerm = 0755

type pullParams struct {
	addr     *url.URL // our sandpiper server
	user     string
	password string
	basePath string
	slice    string // optional (empty means all slices)
	sliceID  uuid.UUID
	debug    bool
}

func getPullParams(c *args.Context) (*pullParams, error) {
	// get sandpiper global params from config file and args
	g, err := GetGlobalParams(c)
	if err != nil {
		return nil, err
	}

	slice := c.String("slice")
	sliceID, _ := uuid.Parse(slice) // ignore error because it might be a slice-name

	return &pullParams{
		addr:     g.addr,
		user:     g.user,
		password: g.password,
		basePath: c.Args().Get(0),
		slice:    slice,
		sliceID:  sliceID,
		debug:    g.debug,
	}, nil
}

type pullCmd struct {
	*pullParams
	api *client.Client
}

// newPullCmd initiates a pull command
func newPullCmd(c *args.Context) (*pullCmd, error) {
	p, err := getPullParams(c)
	if err != nil {
		return nil, err
	}

	// connect to the api server (saving token)
	api, err := client.Login(p.addr, p.user, p.password, p.debug)
	if err != nil {
		return nil, err
	}

	return &pullCmd{pullParams: p, api: api}, nil
}

func (cmd *pullCmd) allSlices() error {
	result, err := cmd.api.ListSlices()
	if err != nil {
		return err
	}
	for _, slice := range result.Slices {
		grain, err := cmd.api.GetLevel1Grain(slice.ID)
		if err != nil {
			return err
		}
		if grain.ID == uuid.Nil {
			fmt.Printf("slice \"%s\" contains no L1 grains\n", slice.Name)
			continue
		}
		if err := saveGrainToFile(cmd.basePath, &slice, grain); err != nil {
			return err
		}
	}
	return nil
}

func (cmd *pullCmd) oneSlice() error {
	var slice *sandpiper.Slice
	var err error

	if cmd.sliceID == uuid.Nil {
		// use provided slice-name to get the slice
		slice, err = cmd.api.SliceByName(cmd.slice)
		if err != nil {
			return err
		}
	} else {
		// use provided slice-id to get the slice
		slice, err = cmd.api.SliceByID(cmd.sliceID)
		if err != nil {
			return err
		}
	}
	grain, err := cmd.api.GetLevel1Grain(slice.ID)
	if err != nil {
		return err
	}
	if err := saveGrainToFile(cmd.basePath, slice, grain); err != nil {
		return err
	}
	return nil
}

// Pull saves file-based grains to the file system
func Pull(c *args.Context) error {

	pull, err := newPullCmd(c)
	if err != nil {
		return err
	}

	if pull.slice == "" {
		return pull.allSlices()
	} else {
		return pull.oneSlice()
	}

}

func saveGrainToFile(basePath string, slice *sandpiper.Slice, grain *sandpiper.Grain) error {

	// todo: change to io.writer (which means we also change the Payload.Decode() code)
	// we probably want to avoid copying payload in memory several times
	/*
		w := bufio.NewWriter(f)
		n4, err := w.WriteString("buffered\n")
		fmt.Printf("wrote %d bytes\n", n4)
		w.Flush()
	*/

	// default filename to grainID if source is empty
	fileName := grain.Source
	if fileName == "" {
		fileName = slice.ID.String() + ".txt"
	}

	// default output to current directory if none provided
	if basePath == "" {
		basePath = "."
	}

	fmt.Printf("Saving: %s/%s/%s ...\n", basePath, slice.Name, fileName)

	folder := fmt.Sprintf("%s/%s", basePath, slice.Name)
	if err := os.MkdirAll(folder, folderPerm); err != nil {
		return fmt.Errorf("unable to create directory \"%s\"", folder)
	}

	f, err := os.Create(folder + "/" + fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	s, err := grain.Payload.Decode(grain.Encoding)
	if err != nil {
		return nil
	}

	if _, err := f.WriteString(s); err != nil {
		return err
	}

	return f.Sync()
}
