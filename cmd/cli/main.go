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

	args "github.com/urfave/cli/v2"

	"autocare.org/sandpiper/pkg/cli"
	"autocare.org/sandpiper/pkg/cli/version"
)

func main() {
	fmt.Println(version.Banner())

	app := args.NewApp()
	app.Name = "sandpiper"
	app.Version = version.Version
	app.Copyright = "Copyright Auto Care Association. All rights reserved."
	app.Usage = "store & retrieve level-1 sandpiper objects"
	app.Flags = cli.GlobalFlags
	app.Commands = cli.Commands
	app.CommandNotFound = cli.CommandNotFound

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
