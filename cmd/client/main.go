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

	"autocare.org/sandpiper/pkg/client/migrations"
	"autocare.org/sandpiper/pkg/client/version"
	"autocare.org/sandpiper/pkg/shared/config"
	"autocare.org/sandpiper/pkg/shared/migrate"
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

	msg := database.Migrate(cfg.DB.URL(), embeddedFiles())
	fmt.Println(msg)

	// todo: poll subscriptions and sync any changes
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
