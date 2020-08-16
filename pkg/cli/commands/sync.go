// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

// Package command implements `sandpiper` commands (add, pull, list, ...)
package command

// sandpiper sync command

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/google/uuid"
	args "github.com/urfave/cli/v2" // conflicts with one of our package names

	"github.com/sandpiper-framework/sandpiper/pkg/shared/client"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/model"
)

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
	// at least one concurrent sync process
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
func (cmd *syncCmd) getServers() (srvs serverList, err error) {
	// lookup is based on command options
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
	fmt.Printf("syncing %s...\n", c.Name)
	return cmd.api.Sync(c)
}

type syncParams struct {
	addr         *url.URL // our sandpiper server
	user         string
	password     string
	partner      string
	partnerID    uuid.UUID
	listOnly     bool
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
		listOnly:     c.Bool("list"),
		noupdate:     c.Bool("noupdate"),
		maxSyncProcs: g.maxSyncProcs,
		debug:        g.debug,
	}, nil
}

// serverList defines a list of active primary servers
type serverList []sandpiper.Company

func (sl serverList) display() {
	if len(sl) == 0 {
		fmt.Println("No active primary servers defined.")
		return
	}
	fmt.Println("SERVER LIST:")
	for _, srv := range sl {
		fmt.Printf("%s: (%s)\n", srv.Name, srv.SyncAddr)
	}
}

// StartSync initiates the sync process on one or all subscriptions
func StartSync(c *args.Context) (result error) {
	var errCount int

	// sync holds params and connected client (to our server)
	sync, err := newSyncCmd(c)
	if err != nil {
		return err
	}

	if sync.api.ServerRole() != sandpiper.SecondaryServer {
		return errors.New("the 'sync' command must be initiated from a secondary server")
	}

	// get list of primary servers to sync (depending on params)
	srvs, err := sync.getServers()
	if err != nil {
		return err
	}

	// display primary servers then exit (--list option)
	if sync.listOnly {
		srvs.display()
		return nil
	}

	// setup a simple waitgroup (using channel as a semaphore and max queue size)
	wg := make(chan int, sync.maxSyncProcs)
	wgAdd := func() { wg <- 1 }
	wgDone := func() { <-wg }
	wgWait := func() {
		for i := 0; i < sync.maxSyncProcs; i++ {
			wgAdd()
		}
	}

	syncFunc := func(srv sandpiper.Company) {
		defer wgDone() // release semaphore
		if err := sync.syncServer(srv); err != nil {
			// show error, but continue (the sync api logs all activity)
			fmt.Printf("syncServer: %v", err)
			result = errors.New("sync completed with errors (see activity logs for details)")
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
