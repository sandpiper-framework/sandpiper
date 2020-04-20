// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package command implements `sandpiper` commands (add, pull, list, ...)
package command

import (
	"fmt"
	"net/url"
	"os"

	"github.com/howeyc/gopass"
	args "github.com/urfave/cli/v2"

	"autocare.org/sandpiper/pkg/shared/config"
	"autocare.org/sandpiper/pkg/sync/client"
)

const (
	// DefaultConfigFile can be overridden by command line options
	DefaultConfigFile = "config.yaml"
)

// GlobalParams holds non-command specific params
type GlobalParams struct {
	addr     *url.URL
	user     string
	password string
	debug    bool
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

	addr, err := getServerAddr(cfg.Command.URL, cfg.Command.Port)
	if err != nil {
		return nil, err
	}

	passwd, err := getPassword(c.String("password"))
	if passwd == "" {
		return nil, fmt.Errorf("password not supplied")
	}

	return &GlobalParams{
		addr:     addr,
		user:     c.String("user"),
		password: passwd,
		debug:    c.Bool("debug"),
	}, nil
}

// Connect to the sandpiper api server (saving token in the client struct)
func Connect(addr *url.URL, user, password string, debug bool) (*client.Client, error) {
	http := client.New(addr, debug)
	if err := http.Login(user, password); err != nil {
		return nil, err
	}
	return http, nil
}

func getPassword(pw string) (string, error) {
	if pw == "" {
		password, err := gopass.GetPasswdPrompt("Password: ", true, os.Stdin, os.Stdout)
		if err != nil {
			return "", err
		}
		pw = string(password)
	}
	return pw, nil
}

func getServerAddr(addr, port string) (*url.URL, error) {
	if port != "" {
		return url.Parse(addr + ":" + port)
	}
	return url.Parse(addr)
}