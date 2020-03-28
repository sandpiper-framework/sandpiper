// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package cli is the sandpiper command line interface
package cli

// command line parser definitions

import (
	"fmt"
	"os"

	args "github.com/urfave/cli/v2" // conflicts with our package name

	"autocare.org/sandpiper/pkg/cli/commands"
)

// sandpiper [global options] command [command options] [arguments]

// global options (also available through env variables)
// 		-u user         # sandpiper user
//		-p password     # sandpiper password
//    -c config       # server access information

// GlobalFlags apply to all the commands
var GlobalFlags = []args.Flag{
	&args.StringFlag{
		Name:     "user",
		Aliases:  []string{"u"},
		Usage:    "api login `name`",
		EnvVars:  []string{"SANDPIPER_USER"},
		Required: true,
	},
	&args.StringFlag{
		Name:    "password",
		Aliases: []string{"p"},
		Usage:   "api password",
		EnvVars: []string{"SANDPIPER_PASSWORD"},
	},
	&args.StringFlag{
		Name:        "config",
		Aliases:     []string{"c"},
		Usage:       "Load configuration from `FILE`",
		EnvVars:     []string{"SANDPIPER_CONFIG"},
		DefaultText: command.DefaultConfigFile,
	},
}

// Commands defines the valid command line sub-commands
var Commands = []*args.Command{
	{
		/*sandpiper add \
		-slice "aap-brake-pads"  \ # slice-name
		-type "aces-file"        \ # grain-type
		-key  "brakes"           \ # grain-key
		-noprompt                \ # don't prompt before over-writing
		"acme_brakes_full_2019-12-12.xml" # file to add (accessed via c.Args().Get(0))
		*/
		Name:      "add",
		Usage:     "add a file-based grain",
		ArgsUsage: "<unzipped-file-to-add>",
		Action:    command.Add,
		Flags: []args.Flag{
			&args.StringFlag{
				Name:     "slice",
				Aliases:  []string{"s"},
				Usage:    "slice name",
				Required: true,
			},
			&args.StringFlag{
				Name:     "type",
				Aliases:  []string{"t"},
				Usage:    "grain-type",
				Required: true,
			},
			&args.StringFlag{
				Name:     "key",
				Aliases:  []string{"k"},
				Usage:    "grain-key",
				Required: true,
			},
			&args.BoolFlag{
				Name:  "noprompt",
				Usage: "do not prompt before over-writing a grain (default is to prompt)",
				Value: false,
			},
		},
	},
	{
		Name:   "pull",
		Usage:  "retrieve file-based grains",
		Action: command.Pull,
		Flags:  []args.Flag{},
	},
	{
		Name:   "list",
		Usage:  "list file-based grains",
		Action: command.List,
		Flags:  []args.Flag{},
	},
}

// CommandNotFound exits program reporting the invalid command
func CommandNotFound(c *args.Context, cmd string) {
	_, _ = fmt.Fprintf(os.Stderr, "%s: '%s' is not a valid command. See '%s --help'.", c.App.Name, cmd, c.App.Name)
	os.Exit(2)
}
