// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package command implements `sandpiper` commands (add, pull, list, ...)
package command

// sandpiper sync command

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/google/uuid"
	args "github.com/urfave/cli/v2" // conflicts with one of our package names

	"sandpiper/pkg/shared/client"
	"sandpiper/pkg/shared/model"
)

// serverList defines a list of active primary servers
type serverList []sandpiper.Company

type syncParams struct {
	addr         *url.URL // our sandpiper server
	user         string
	password     string
	partner      string
	partnerID    uuid.UUID
	noupdate     bool
	maxSyncProcs int
	debug        bool
}

func getSyncParams(c *args.Context) (*syncParams, error) {
	// get global params from config file and args
	g, err := GetGlobalParams(c)
	if err != nil {
		return nil, err
	}
	// save optional "partner" param, and determine its type
	p := c.String("partner")
	id, _ := uuid.Parse(p) // valid id means companyID, otherwise company-name

	return &syncParams{
		addr:         g.addr,
		user:         g.user,
		password:     g.password,
		partner:      p,
		partnerID:    id,
		noupdate:     c.Bool("noupdate"),
		maxSyncProcs: g.maxSyncProcs,
		debug:        g.debug,
	}, nil
}

type syncCmd struct {
	*syncParams
	api *client.Client
}

// newSyncCmd initiates a sync command
func newSyncCmd(c *args.Context) (*syncCmd, error) {
	p, err := getSyncParams(c)
	if err != nil {
		return nil, err
	}
	// make sure we allow at least one concurrent sync process
	if p.maxSyncProcs <= 0 {
		p.maxSyncProcs = 1
	}
	// connect to our api server (saving token)
	api, err := client.Login(p.addr, p.user, p.password, p.debug)
	if err != nil {
		return nil, err
	}
	return &syncCmd{syncParams: p, api: api}, nil
}

// getServers returns list of servers to sync with
func (cmd *syncCmd) getServers() (serverList, error) {
	var srvs serverList
	var err error

	switch {
	case cmd.partner == "":
		// retrieve all syncable servers
		srvs, err = cmd.api.ActiveServers(uuid.Nil, "")
	case cmd.partnerID != uuid.Nil:
		// use provided company_id to get a syncable server
		srvs, err = cmd.api.ActiveServers(cmd.partnerID, "")
	default:
		// use provided partner-name to get the company's server info
		srvs, err = cmd.api.ActiveServers(uuid.Nil, cmd.partner)
	}

	if cmd.debug && err != nil && srvs != nil {
		// display the list of sync servers to console
		srvs.display()
	}

	return srvs, err
}

// syncServer performs the actual sync on a server
func (cmd *syncCmd) syncServer(c sandpiper.Company) error {
	if cmd.debug {
		fmt.Printf("syncing %s...", c.Name)
	}
	return cmd.api.Sync(c)
}

// StartSync initiates the sync process on one or all subscriptions
func StartSync(c *args.Context) error {
	var result error
	var errCount int

	// sync holds params and connected client (to our server)
	sync, err := newSyncCmd(c)
	if err != nil {
		return err
	}

	// get list of primary servers to sync (depending on params)
	srvs, err := sync.getServers()
	if err != nil {
		return err
	}

	// setup a simple waitgroup (using channel as a semaphore and max queue size)
	wg := make(chan struct{}, sync.maxSyncProcs)
	wgAdd := func() { wg <- struct{}{} }
	wgDone := func() { <-wg }
	wgWait := func() {
		for i := 0; i < sync.maxSyncProcs; i++ {
			wgAdd()
		}
	}

	syncFunc := func(srv sandpiper.Company) {
		defer wgDone() // release semaphore
		if err := sync.syncServer(srv); err != nil {
			// log error, but continue
			if sync.debug {
				fmt.Printf("%v", err)
			}
			// todo: log error to activity
			result = errors.New("sync completed with errors")
			errCount++
		}
	}

	// sync each server in a separate go routine
	for _, srv := range srvs {
		wgAdd() // acquire a semaphore
		go syncFunc(srv)
	}
	wgWait() // wait for all to finish

	fmt.Printf("Successful: %d, Errors: %d\n", len(srvs)-errCount, errCount)
	return result
}

func (sl serverList) display() {
	fmt.Println("SERVER LIST:")
	for _, srv := range sl {
		fmt.Printf("%s: (%s)\n", srv.Name, srv.SyncAddr)
	}
}
