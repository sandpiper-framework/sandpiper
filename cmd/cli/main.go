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
	"strconv"
	"time"

	args "github.com/urfave/cli/v2" // conflicts with our package name

	"sandpiper/pkg/cli"
	"sandpiper/pkg/cli/version"
)

func main() {
	fmt.Println(version.Banner())

	y := strconv.Itoa(time.Now().Year())

	app := &args.App{
		Name:     "sandpiper",
		Version:  version.Version,
		Compiled: time.Now(),
		Authors: []*args.Author{
			{
				Name:  "Doug Winsby",
				Email: "dougw@winsbygroup.com",
			},
		},
		Copyright:       "Copyright 2019-" + y + " The Sandpiper Authors. All rights reserved.",
		HelpName:        "sandpiper",
		Usage:           "Store, extract, list and sync \"level-1\" (file-based) sandpiper objects",
		Flags:           cli.GlobalFlags,
		Commands:        cli.Commands,
		CommandNotFound: cli.CommandNotFound,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
