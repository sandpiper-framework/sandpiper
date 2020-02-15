// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package main is the entry point for the `sandpiper` command.
// It opens the database (using a config file), checks the database-version and executes the provided command.
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	args "github.com/urfave/cli/v2"

	"autocare.org/sandpiper/pkg/cli"
	"autocare.org/sandpiper/pkg/cli/version"
)

func main() {
	fmt.Println(version.Banner())

	app := &args.App{
		Name: "sandpiper",
		Version: version.Version,
		Compiled: time.Now(),
		Authors: []*args.Author{
			&args.Author{
				Name:  "Doug Winsby",
				Email: "dougw@winsbygroup.com",
			},
		},
		Copyright: "Copyright Auto Care Association. All rights reserved.",
		HelpName: "sandpiper",
		Usage: "Store & retrieve \"level-1\" (file-based) sandpiper objects",
		Flags: cli.GlobalFlags,
		Commands: cli.Commands,
		CommandNotFound: cli.CommandNotFound,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
