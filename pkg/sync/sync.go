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
		--slice "aap-brake-pads"  \ # argument is a slice name
		--noupdate                \ # perform the sync without actually changing anything
		*/
		Name:      "run",
		Usage:     "begin a sync process on all active subscriptions",
		ArgsUsage: " ",
		Action:    command.Run,
		Flags: []args.Flag{
			&args.StringFlag{
				Name:     "slice",
				Aliases:  []string{"s"},
				Usage:    "either a slice_id (uuid) or slice_name (case-insensitive)",
				Required: false,
			},
			&args.BoolFlag{
				Name:  "noupdate",
				Usage: "perform the sync without actually changing anything",
			},
		},
	},
	{
		/* sync show \
		--slice "aap-slice" \ # required slice_id or slice_name
		arg  # either slice_id or slice_name (if empty, list all slices)
		*/
		Name:      "show",
		Usage:     "list slices (if no slice provided) or file-based grains by slice_id or slice_name",
		ArgsUsage: " ", // don't show that we accept arguments
		Action:    command.Show,
		Flags: []args.Flag{
			&args.StringFlag{
				Name:     "slice",
				Aliases:  []string{"s"},
				Usage:    "either a slice_id (uuid) or slice_name (case-insensitive)",
				Required: false,
			},
			&args.BoolFlag{
				Name:  "full",
				Usage: "provide full listings",
			},
		},
	},
}

// CommandNotFound exits program reporting the invalid command
func CommandNotFound(c *args.Context, cmd string) {
	_, _ = fmt.Fprintf(os.Stderr, "%s: '%s' is not a valid command. See '%s --help'.", c.App.Name, cmd, c.App.Name)
	os.Exit(2)
}
