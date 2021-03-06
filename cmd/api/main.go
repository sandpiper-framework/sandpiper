// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

// Package main is the entry point for the sandpiper primary server.
// It checks the database-version and uses a config file to launch the server properly.
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/sandpiper-framework/sandpiper/pkg/api"
	"github.com/sandpiper-framework/sandpiper/pkg/api/version"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/config"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/database"
)

const (
	debugModeMsg = "** RUNNING IN DEBUG (NON-PRODUCTION) MODE **"
)

func main() {
	fmt.Println(version.Banner())

	cfgPath := flag.String("config", "api-config.yaml", "Path to config file (default is api-config.yaml)")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatal("ERROR: ", err)
	}

	// Update the database if necessary from code (as defined in "shared/database/schema.go")
	msg, err := database.Migrate(cfg.DB.DSN())
	if err != nil {
		log.Fatal("ERROR: ", err)
	}
	fmt.Printf("Database: \"%s\"\n%s\n", cfg.DB.Database, msg)

	if cfg.Server.Debug {
		fmt.Printf("\n%s\n%s\n", debugModeMsg, cfg.DB.SafeDSN())
	}

	err = api.Start(cfg)
	if err != nil {
		log.Fatal("ERROR: ", err)
	}
}
