// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package main is the entry point for the `sandpiper` command.
// It uses the published api to perform operations against a primary or secondary server.
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	args "github.com/urfave/cli/v2" // conflicts with our package name

	"autocare.org/sandpiper/pkg/sync"
	"autocare.org/sandpiper/pkg/sync/version"
)

func main() {
	fmt.Println(version.Banner())

	app := &args.App{
		Name:     "sync",
		Version:  version.Version,
		Compiled: time.Now(),
		Authors: []*args.Author{
			{
				Name:  "Doug Winsby",
				Email: "dougw@winsbygroup.com",
			},
		},
		Copyright:       "Copyright The Sandpiper Authors. All rights reserved.",
		HelpName:        "sync",
		Usage:           "perform a Sandpiper `sync` operation with all trading partners",
		Flags:           sync.GlobalFlags,
		Commands:        sync.Commands,
		CommandNotFound: sync.CommandNotFound,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
