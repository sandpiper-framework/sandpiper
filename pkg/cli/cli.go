// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

// Package cli is the sandpiper command line interface
package cli

// command line parser definitions

import (
	"fmt"
	"os"

	args "github.com/urfave/cli/v2" // conflicts with our package name

	"github.com/sandpiper-framework/sandpiper/pkg/cli/commands"
)

// sandpiper [global options] command [command options] [arguments]

// global options (also available through env variables)
// 		-u user         # sandpiper user
//		-p password     # sandpiper password
//    -c config       # server access information

// GlobalFlags apply to all the commands
var GlobalFlags = []args.Flag{
	&args.StringFlag{
		Name:    "user",
		Aliases: []string{"u"},
		Usage:   "server login `name`",
		EnvVars: []string{"SANDPIPER_USER"},
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
		/* sandpiper add \
		   --slice "aap-brake-pads"  \ # argument is a slice name
		   --noprompt                \ # don't prompt before over-writing
		   acme_brakes_full_2019-12-12.xml # file to add (accessed via c.Args().Get(0))
		*/
		Name:      "add",
		Usage:     "add a file-based grain from a local file",
		ArgsUsage: "<unzipped-file-to-add>",
		Action:    command.Add,
		Flags: []args.Flag{
			&args.StringFlag{
				Name:     "slice",
				Aliases:  []string{"s"},
				Usage:    "either slice_id (uuid) or slice_name (case-insensitive)",
				Required: true,
			},
			&args.BoolFlag{
				Name:  "noprompt",
				Usage: "do not prompt before over-writing a grain (default is to prompt)",
			},
		},
	},
	{
		/* sandpiper pull \
		   --slice "aap-slice" \ # required slice_id or slice_name
		   --dir                 # required output directory
		*/
		Name:      "pull",
		Usage:     "save file-based grains to the file system",
		ArgsUsage: "<output-directory>",
		Action:    command.Pull,
		Flags: []args.Flag{
			&args.StringFlag{
				Name:     "slice",
				Aliases:  []string{"s"},
				Usage:    "either a slice_id (uuid) or slice_name (case-insensitive)",
				Required: false,
			},
			&args.StringFlag{
				Name:     "dir",
				Aliases:  []string{"d"},
				Usage:    "optional root of output directory",
				Required: false,
			},
		},
	},
	{
		/* sandpiper list \
		   --slice "aap-slice"  # slice_id or slice_name
		*/
		Name:      "list",
		Usage:     "list slices (if no slice provided) or file-based grains by slice_id or slice_name",
		ArgsUsage: " ", // don't show that we accept arguments
		Action:    command.List,
		Flags: []args.Flag{
			&args.StringFlag{
				Name:     "slice",
				Aliases:  []string{"s"},
				Usage:    "either a slice_id (uuid) or slice_name (case-insensitive)",
				Required: false,
			},
			&args.BoolFlag{
				Name:     "full",
				Usage:    "provide full listings",
				Required: false,
			},
		},
	},
	{
		/* sandpiper sync \
		   --company "acme-brakes"  \ # an optional company name (case-insensitive) or company_id
		   --list                   \ # show active servers without performing a sync
		   --noupdate                 # perform the sync without actually changing anything
		*/
		Name:      "sync",
		Usage:     "Start the sync process on active subscriptions",
		ArgsUsage: " ", // no arguments
		Action:    command.StartSync,
		Flags: []args.Flag{
			&args.StringFlag{
				Name:     "partner",
				Aliases:  []string{"p"},
				Usage:    "limit to company name (case-insensitive) or company_id",
				Required: false,
			},
			&args.BoolFlag{
				Name:     "list",
				Aliases:  []string{"l"},
				Usage:    "Display a list of sync servers (without performing a sync)",
				Required: false,
			},
			&args.BoolFlag{
				Name:     "noupdate",
				Usage:    "Perform the sync without actually changing anything locally",
				Required: false,
			},
		},
	},
	{
		/* sandpiper init
		 */
		Name:      "init",
		Usage:     "initialize a sandpiper primary or secondary database",
		ArgsUsage: " ", // don't show that we accept arguments
		Action:    command.Init,
		Flags: []args.Flag{
			&args.StringFlag{
				Name:  "id",
				Usage: "assign this server-id (for testing only)",
			},
			&args.BoolFlag{
				Name:     "debug",
				Usage:    "show debug messages during the init process",
				Required: false,
			},
		},
	},
	{
		/* sandpiper secrets
		 */
		Name:      "secrets",
		Usage:     "generate new random secrets for env vars and api-config.yaml file",
		ArgsUsage: " ", // don't show that we accept arguments
		Action:    command.Secrets,
	},
}

// CommandNotFound exits program reporting the invalid command
func CommandNotFound(c *args.Context, cmd string) {
	_, _ = fmt.Fprintf(os.Stderr, "%s: '%s' is not a valid command. See '%s --help'.", c.App.Name, cmd, c.App.Name)
	os.Exit(2)
}
