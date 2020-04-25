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

	"autocare.org/sandpiper/pkg/cli/client"
	"autocare.org/sandpiper/pkg/shared/model"
)

// subArray defines a list of subscription
type subsArray []sandpiper.Subscription

type syncParams struct {
	addr         *url.URL // our sandpiper server
	user         string
	password     string
	subscription string // required
	companyID    uuid.UUID
	noupdate     bool
	debug        bool
}

func getSyncParams(c *args.Context) (*syncParams, error) {
	// get global params from config file and args
	g, err := GetGlobalParams(c)
	if err != nil {
		return nil, err
	}
	// save optional "subscription" param, and determine its type
	sub := c.String("subscription")
	id, _ := uuid.Parse(sub) // valid id means companyID, otherwise subscription-name

	return &syncParams{
		addr:         g.addr,
		user:         g.user,
		password:     g.password,
		subscription: sub,
		companyID:    id,
		noupdate:     c.Bool("noupdate"),
		debug:        g.debug,
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
	api, err := Connect(p.addr, p.user, p.password, p.debug)
	if err != nil {
		return nil, err
	}
	return &syncCmd{syncParams: p, api: api}, nil
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
		subs, err = cmd.api.SubsByName(cmd.subscription)
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

// syncServer performs the actual sync on a server
func (subs subsArray) syncServer(addr string) error {
	// open websocket with primary server

	//
	for _, sub := range subs {
		slice := sub.Slice
		if !slice.AllowSync {
			// just log that it is locked
		}

	}
	return nil
}

// StartSync initiates the sync process on one or all subscriptions
func StartSync(c *args.Context) error {
	sync, err := newSyncCmd(c)
	if err != nil {
		return err
	}
	subs, err := sync.subscriptions()
	if err != nil {
		return err
	}
	return subs.sync(sync.debug)
}
