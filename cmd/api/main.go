// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package main is the entry point for the sandpiper primary server.
// It checks the database-version and uses a config file to launch the server properly.
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4/source/go_bindata"

	"sandpiper/pkg/api"
	"sandpiper/pkg/api/migrations"
	"sandpiper/pkg/api/version"
	"sandpiper/pkg/shared/config"
	"sandpiper/pkg/shared/migrate"
)

const (
	debugModeMsg = "** RUNNING IN DEBUG (NON-PRODUCTION) MODE **"
)

func main() {
	fmt.Println(version.Banner())

	cfgPath := flag.String("config", "config.yaml", "Path to config file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatal("ERROR: ", err)
	}

	// Update the database if necessary (from bindata embedded files)
	msg := database.Migrate(cfg.DB.URL(), embeddedFiles())
	fmt.Printf("Database: \"%s\"\n%s\n", cfg.DB.Database, msg)

	if cfg.Server.Debug {
		fmt.Printf("\n%s\n%s\n", debugModeMsg, cfg.DB.SafeURL())
	}

	err = api.Start(cfg)
	if err != nil {
		log.Fatal("ERROR: ", err)
	}
}

// embeddedFiles returns a pointer to the structure that manages access to embedded database migration files.
// It uses an "import" specific to the pkg we are building (so this function must be local for each executable).
func embeddedFiles() *bindata.AssetSource {
	r := bindata.Resource(migrations.AssetNames(),
		func(name string) ([]byte, error) {
			return migrations.Asset(name)
		})
	return r
}
