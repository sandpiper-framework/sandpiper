package database

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Migrate applies database migrations necessary to bring the database-version up to date.
// Report any problems with the standard logger (not an api logger)
func Migrate(psn string, migrationDir string) string {

	// Read migrations from directory and connect to database.
	m, err := migrate.New("file://"+migrationDir, psn)
	if err != nil {
		log.Fatal(err)
	}

	v1, v2 := applyMigrations(m)

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
	return fmt.Sprintf("DB Version: %d (no migrations required)\n", v2)
}