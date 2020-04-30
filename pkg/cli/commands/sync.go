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

	"autocare.org/sandpiper/pkg/shared/client"
	"autocare.org/sandpiper/pkg/shared/model"
)

// serverList defines a list of active primary servers
type serverList []sandpiper.Company

// subArray defines a list of subscription
//type subsArray []sandpiper.Subscription

type syncParams struct {
	addr      *url.URL // our sandpiper server
	user      string
	password  string
	partner   string
	partnerID uuid.UUID
	noupdate  bool
	debug     bool
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
		addr:      g.addr,
		user:      g.user,
		password:  g.password,
		partner:   p,
		partnerID: id,
		noupdate:  c.Bool("noupdate"),
		debug:     g.debug,
	}, nil
}

type syncCmd struct {
	*syncParams
	api *client.Client
}

// newSyncCmd initiates a run command
func newSyncCmd(c *args.Context) (*syncCmd, error) {
	p, err := getSyncParams(c)
	if err != nil {
		return nil, err
	}
	// connect to the api server (saving token)
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

	sync, err := newSyncCmd(c)
	if err != nil {
		return err
	}

	// get list of primary servers to sync (depending on run params)
	srvs, err := sync.getServers()
	if err != nil {
		return err
	}

	// todo: make each sync a go routine (https://gobyexample.com/worker-pools)
	/* COMMENT:

	I think most of the time you probably don't really need a pool. But maybe you do want to limit the number of goroutines
	running at once.

	E.g. if each goroutine is CPU intensive and you want to echo progress on the tasks in a reasonable manner, rather than
	having all of them complete around the same time.

	An easy way to limit how many goroutines are running at once is to use a channel as a semaphore. E.g. something like

	sema := make(chan struct{}, numCores)
	and then in your worker goroutines:

	sema <- struct{}{} // blocks if more than numCores threads doing  work
	// ... do work ... //
	<-sema

	That'll ensure you don't have more than numCores goroutines running at once. IMO much easier than a pool, but you
	still have the nice property that you're not making tiny progress on many tasks at once, rather just completing tasks
	as quickly as possible.

	If you really do want a pool, then using a channel to send in the work seems reasonable. If you want to get a response,
	you can actually include a response channel in your message on the input channel.

	*/

	// sync each server
	for _, srv := range srvs {
		if err := sync.syncServer(srv); err != nil {
			// log error, but continue
			if sync.debug {
				fmt.Printf("%v", err)
			}
			// todo: log error to activity
			result = errors.New("sync completed with errors")
		}
	}

	return result
}

//********************** dead code we may want ****************

/*

// updateSubs asks a server for all of our subscriptions so we can keep them updated
func (cmd *syncCmd) updateSubs() error {

	subs, err := cmd.api.AllSubs()
	if err != nil {
		return err
	}
	for _, sub := range subs {
		slice := sub.Slice
		if !slice.AllowSync {
			// just log that it is locked
		}

	}

	return nil
}

// subscriptions returns an array of subscriptions to sync
func (cmd *syncCmd) subscriptions() (subsArray, error) {
	var (
		subs subsArray
		err  error
	)

	switch {
	case cmd.subscription == "":
		// retrieve all subscriptions
		subs, err = cmd.api.AllSubs()
	case cmd.companyID != uuid.Nil:
		// use provided company_id to get the subscription(s)
		subs, err = cmd.api.SubsByCompany(cmd.companyID)
	default:
		// use provided subscription-name to get the subscription
		sub, err := cmd.api.SubByName(cmd.subscription)
		if err != nil {
			return nil, err
		}
		subs = append(subs, *sub)
	}
	return subs, err
}

// sync organizes the work to do and calls the sync routine for each subscription
func (subs subsArray) sync(debugFlag bool) error {
	var result error

	// organize active subs by syncAddr using a "multimap" [syncAddr: subs]
	work := make(map[string]subsArray)
	for _, sub := range subs {
		if sub.Active && sub.Company.Active {
			// add to the list of subscriptions to sync for this sync_addr
			addr := sub.Company.SyncAddr
			work[addr] = append(work[addr], sub)
		}
	}
	// sync each syncAddr, sync subscriptions
	for addr, subs := range work {
		if err := subs.syncServer(addr); err != nil {
			// log error, but continue
			if debugFlag {
				fmt.Printf("%v", err)
			}
			result = errors.New("sync completed with errors")
		}
	}
	return result
}

*/
