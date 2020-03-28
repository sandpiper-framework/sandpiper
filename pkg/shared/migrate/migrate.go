// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package database

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/go_bindata"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)

// Migrate applies database migrations necessary to bring the database-version up to date.
// These migrations are embedded in a source file by `go-bindata`. Reports problems using
// the standard logger (not an api logger)
func Migrate(psn string, bin *bindata.AssetSource) string {
	src, err := bindata.WithInstance(bin)
	if err != nil {
		log.Fatal(err)
	}
	m, err := migrate.NewWithSourceInstance("go-bindata", src, psn)
	if err != nil {
		log.Fatal(err)
	}

	v1, v2 := applyMigrations(m)
	if v1 > v2 {
		log.Fatalf("software version mismatch: db version (%d) is more current than program expects (%d)... update the software", v1, v2)
	}

	_, err = m.Close()
	if err != nil {
		log.Fatal(err)
	}

	return migrationMessage(v1, v2)
}

// applyMigrations looks for all the "up" migrations after the current database-version (if any) and
// runs those sql files in numerical order. Return the before and after version numbers.
func applyMigrations(m *migrate.Migrate) (uint, uint) {

	// get version before migration
	v1 := getDatabaseVersion(m)

	// Migrate all the way up ...
	err := m.Up()
	if err != nil {
		// for some reason, "no change" is an error
		if err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}

	// get version after migration
	v2 := getDatabaseVersion(m)

	return v1, v2
}

// getDatabaseVersion returns the current database version
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

func migrationMessage(v1, v2 uint) string {
	if v1 != v2 {
		return fmt.Sprintf("DB Version: %d (migrated from %d to %d)\n", v2, v1, v2)
	}
	return fmt.Sprintf("DB Version: %d\n", v1)
}
