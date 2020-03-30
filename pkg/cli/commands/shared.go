// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package command implements `sandpiper` commands (add, pull, list, ...)
package command

import (
	"bufio"
	"fmt"
	args "github.com/urfave/cli/v2"
	"net/url"
	"os"
	"strings"

	"autocare.org/sandpiper/pkg/cli/client"
	"autocare.org/sandpiper/pkg/shared/config"
)

const (
	// DefaultConfigFile can be overridden by command line options
	DefaultConfigFile = "cli.config.yaml"
)

// GlobalParams holds non-command specific params
type GlobalParams struct {
	addr     *url.URL
	user     string
	password string
}

// GetGlobalParams parses global parameters from command line
func GetGlobalParams(c *args.Context) (*GlobalParams, error) {

	cfgPath := c.String("config")
	if cfgPath == "" {
		cfgPath = DefaultConfigFile
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return nil, err
	}

	addr, err := url.Parse(cfg.Command.URL)
	if err != nil {
		return nil, err
	}

	return &GlobalParams{
		addr:     addr,
		user:     c.String("user"),
		password: c.String("password"),
	}, nil
}

// Connect to the sandpiper api server (saving token in the client struct)
func Connect(addr *url.URL, user, password string) (*client.Client, error) {
	http := client.New(addr)
	if err := http.Login(user, password); err != nil {
		return nil, err
	}
	return http, nil
}

// AllowOverwrite prompts the user for permission to overwrite something
func AllowOverwrite() bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Overwrite (y/n)? ")
	ans, _ := reader.ReadString('\n')
	return strings.ToLower(ans) == "y"
}
