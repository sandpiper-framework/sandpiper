// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package sync is the sandpiper sync command
package sync

// command line parser definitions

import (
	"fmt"
	"os"

	args "github.com/urfave/cli/v2" // conflicts with our package name

	"autocare.org/sandpiper/pkg/sync/commands"
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
		Usage:    "server login `name`",
		EnvVars:  []string{"SANDPIPER_USER"},
		Required: true,
	},
	&args.StringFlag{
		Name:    "password",
		Aliases: []string{"p"},
		Usage:   "user password",
		EnvVars: []string{"SANDPIPER_PASSWORD"},
	},
	&args.StringFlag{
		Name:        "config",
		Aliases:     []string{"c"},
		Usage:       "Load configuration from `FILE`",
		EnvVars:     []string{"SANDPIPER_CONFIG"},
		DefaultText: command.DefaultConfigFile,
	},
	&args.BoolFlag{
		Name:    "debug",
		Aliases: []string{"d"},
		Usage:   "Show debug information to stdout",
	},
}

// Commands defines the valid command line sub-commands
var Commands = []*args.Command{
	{
		/* sync run \
		--slice "aap-brake-pads"  \ # an optional slice-id or slice-name
		--noupdate                  # perform the sync without actually changing anything
		*/
		Name:      "run",
		Usage:     "Start the sync process on all active subscriptions",
		ArgsUsage: " ", // no arguments
		Action:    command.Run,
		Flags: []args.Flag{
			&args.StringFlag{
				Name:     "slice",
				Aliases:  []string{"s"},
				Usage:    "Either a slice_id (uuid) or slice_name (case-insensitive)",
				Required: false,
			},
			&args.BoolFlag{
				Name:  "noupdate",
				Usage: "Perform the sync without actually changing anything",
			},
		},
	},
	{
		/* sync show \
		--slice "aap-slice" \ # slice_id or slice_name
		--full                # detailed information
		*/
		Name:      "show",
		Usage:     "Display sync information for all slices (if none provided) or a single slice by slice_id or slice_name",
		ArgsUsage: " ", // don't show that we accept arguments
		Action:    command.Show,
		Flags: []args.Flag{
			&args.StringFlag{
				Name:     "slice",
				Aliases:  []string{"s"},
				Usage:    "Either a slice_id (uuid) or slice_name (case-insensitive)",
				Required: false,
			},
			&args.BoolFlag{
				Name:  "full",
				Usage: "Provide detailed information",
			},
		},
	},
}

// CommandNotFound exits program reporting the invalid command
func CommandNotFound(c *args.Context, cmd string) {
	_, _ = fmt.Fprintf(os.Stderr, "%s: '%s' is not a valid command. See '%s --help'.", c.App.Name, cmd, c.App.Name)
	os.Exit(2)
}
