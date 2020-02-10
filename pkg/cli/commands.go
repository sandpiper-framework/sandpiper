package cli

import (
	"fmt"
	"os"

	args "github.com/urfave/cli"

	"autocare.org/sandpiper/pkg/cli/commands"
)

var GlobalFlags = []args.Flag{
	
}

// Commands defines the valid command line sub-commands
var Commands = []args.Command{
	{
		Name:        "add",
		Usage:       "add a file-based grain",
		Action:      command.Add,
		Flags:       []args.Flag{},
	},
	{
		Name:        "pull",
		Usage:       "retrieve file-based grains",
		Action:      command.Pull,
		Flags:       []args.Flag{},
	},
	{
		Name:        "list",
		Usage:       "list file-based grains",
		Action:      command.List,
		Flags:       []args.Flag{},
	},
}

func CommandNotFound(c *args.Context, cmd string) {
	_, _ = fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, cmd, c.App.Name, c.App.Name)
	os.Exit(2)
}
