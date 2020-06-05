// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

// Package command implements `sandpiper` commands (add, pull, list, ...)
package command

import (
	"bufio"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/howeyc/gopass"
	args "github.com/urfave/cli/v2"

	"sandpiper/pkg/shared/config"
)

const (
	// DefaultConfigFile can be overridden by command line options
	DefaultConfigFile = "cli-config.yaml"

	// L1Encoding for level-1 grains using `sandpiper add`
	L1Encoding = "z64"
)

// GlobalParams holds non-command specific params
type GlobalParams struct {
	addr         *url.URL
	user         string
	password     string
	maxSyncProcs int
	debug        bool
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

	if cfg.Command == nil {
		return nil, errors.New("config file must have a \"command\" section")
	}

	addr, err := getServerAddr(cfg.Command.URL, cfg.Command.Port)
	if err != nil {
		return nil, err
	}

	user := c.String("user")
	if user == "" {
		return nil, errors.New("user not supplied")
	}

	passwd := getPassword(c.String("password"))
	if passwd == "" {
		return nil, errors.New("password not supplied")
	}

	return &GlobalParams{
		addr:         addr,
		user:         c.String("user"),
		password:     passwd,
		maxSyncProcs: cfg.Command.MaxSyncProcs,
		debug:        c.Bool("debug"),
	}, nil
}

// AllowOverwrite prompts the user for permission to overwrite something
func AllowOverwrite() bool {
	ans := Prompt("Overwrite (y/n)? ", "y")
	return strings.ToLower(ans) == "y"
}

// Prompt for an answer given a question
func Prompt(question, defaultAns string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(question)
	ans, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}
	ans = strings.TrimSpace(ans)
	if ans == "" {
		ans = defaultAns
	}
	return ans
}

// GetPassword accepts console input masking typed characters
func GetPassword(prompt string) string {
	pw, err := gopass.GetPasswdPrompt(prompt, true, os.Stdin, os.Stdout)
	if err != nil {
		return ""
	}
	return string(pw)
}

func getPassword(pw string) string {
	if pw == "" {
		pw = GetPassword("Password: ")
	}
	return pw
}

func getServerAddr(addr, port string) (*url.URL, error) {
	if port != "" {
		return url.Parse(addr + ":" + port)
	}
	return url.Parse(addr)
}
