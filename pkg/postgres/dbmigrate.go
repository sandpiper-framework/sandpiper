package postgres

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// MigrateDatabase applies migrations necessary to bring the database up to date.
func MigrateDatabase(psn string, migrationDir string) (msg string) {

	// Read migrations from directory and connect to database.
	m, err := migrate.New("file://"+migrationDir, psn)
	if err != nil {
		log.Fatal(err)
	}

	oldVer := getDatabaseVersion(m)

	// Migrate all the way up ...
	err = m.Up()
	if err != nil {
		if fmt.Sprintf("%v", err) != "no change" { // todo: a better comparison?!
			log.Fatal(err)
		}
	}

	newVer := getDatabaseVersion(m)
	if newVer != oldVer {
		msg = fmt.Sprintf("DB Version: %d (migrations applied from %d)\n", oldVer, newVer)
	} else {
		msg = fmt.Sprintf("DB Version: %d (no migrations required)\n", newVer)
	}

	_, _ = m.Close()
	return
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