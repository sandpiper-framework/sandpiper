// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package main is the entry point for the sandpiper client.
// It checks the database-version and uses a config file to launch the server properly.
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4/source/go_bindata"

	"autocare.org/sandpiper/internal/config"
	"autocare.org/sandpiper/internal/database"
	"autocare.org/sandpiper/pkg/client/migrations"
	"autocare.org/sandpiper/pkg/client/version"
)

func main() {
	fmt.Println(version.Banner())

	cfgPath := flag.String("p", "./client.config.yaml", "Path to config file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatal("ERROR: ", err)
	}

	// Update the database if necessary
	bin := assetResource()
	msg := database.Migrate(cfg.DB.PSN, bin)
	fmt.Println(msg)

	// todo: do stuff
}

// assetResource returns a pointer to the structure that manages access to embedded migration files.
// It uses migrations specific to the pkg we are building (based on the correct import).
func assetResource() *bindata.AssetSource {
	r := bindata.Resource(migrations.AssetNames(),
		func(name string) ([]byte, error) {
			return migrations.Asset(name)
		})
	return r
}