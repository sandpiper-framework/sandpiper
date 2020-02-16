// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package command implements `sandpiper` commands (add, pull, list, ...)
package command

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"

	"autocare.org/sandpiper/pkg/cli/client"
)

const (
	// DefaultConfigFile can be overridden by command line options
	DefaultConfigFile = "cli.config.yaml"
)

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
