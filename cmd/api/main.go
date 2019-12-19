// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package main is the entry point for the sandpiper api server.
// It checks the database-version and uses a config file to launch the server properly
package main

import (
	"flag"
	"fmt"
	"log"

	"autocare.org/sandpiper/internal/config"
	"autocare.org/sandpiper/internal/database"
	"autocare.org/sandpiper/internal/version"
	"autocare.org/sandpiper/pkg/api"
)

func main() {
	fmt.Println(version.Banner())

	cfgPath := flag.String("p", "./sandpiper.config.yaml", "Path to config file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatal("ERROR: ", err)
	}

	// Update the database if necessary
	msg := database.Migrate(cfg.DB.PSN, cfg.DB.MigrateDir)
	fmt.Println(msg)

	err = api.Start(cfg)
	if err != nil {
		panic(err.Error())
	}
}
