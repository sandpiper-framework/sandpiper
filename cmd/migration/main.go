package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"autocare.org/sandpiper/pkg/config"
)

const migrationFolder = "./migrations"

func main() {
	cfgPath := flag.String("p", "../api/sandpiper.config.yaml", "Path to config file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("%s: %v", "config not loaded", err)
	}

	// Read migrations from directory and connect to database.
	m, err := migrate.New("file://"+migrationFolder, cfg.DB.PSN)
	if err != nil {
		log.Fatal(err)
	}

	oldVer := getDatabaseVersion(m)

	// Migrate all the way up ...
	err = m.Up()
	if err != nil {
		if fmt.Sprintf("%v", err) != "no change" { // todo: a better way?!
			log.Fatal(err)
		}
	}

	newVer := getDatabaseVersion(m)
	if newVer != oldVer {
		log.Printf("Database migrations applied (%d to %d)", oldVer, newVer)
	} else {
		log.Printf("db version: %d (no migrations required)", newVer)
	}

	_, _ = m.Close()
}

func getDatabaseVersion(m *migrate.Migrate) uint {
	ver, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Fatal(err)
	}
	if dirty {
		log.Fatal("Previous failed migration found, contact database administrator.")
	}
	return ver
}
