// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

// Package database creates a pooled connection to the database. We use a
// lightweight ORM (with deep support for postgresql). This ORM only supports
// postgresql. We might consider switching if require support for other dbms.
package database

import (
	"context"
	"crypto/tls"
	"errors"
	"log"
	"os"
	"time"

	"github.com/go-pg/pg/v9"
	// DB adapter
	_ "github.com/lib/pq"

	"sandpiper/pkg/shared/config"
	"sandpiper/pkg/shared/model"
)

type dbLogger struct{}

// DB is a wrapper around pg.DB so we can add functionality
type DB struct {
	*pg.DB
	Settings *sandpiper.Setting // database settings
}

// BeforeQuery is an unused stub at this time.
func (d dbLogger) BeforeQuery(ctx context.Context, _ *pg.QueryEvent) (context.Context, error) {
	return ctx, nil
}

// AfterQuery is a callback hook allowing us to log the query if in debug mode.
func (d dbLogger) AfterQuery(_ context.Context, q *pg.QueryEvent) error {
	log.Printf(q.FormattedQuery())
	return nil
}

// New creates new database connection to a postgres database with optional query logging
func New(c *config.Database, timeout int, enableLog bool) (*DB, error) {

	// wrap db connection in our struct
	db := &DB{DB: pg.Connect(ConnectOptions(c))}

	// test connectivity
	if _, err := db.Exec("SELECT 1"); err != nil {
		return nil, err
	}

	// save any database settings in our db object
	if err := db.settings(); err != nil {
		return nil, err
	}

	// make configuration settings
	if timeout > 0 {
		db.WithTimeout(time.Second * time.Duration(timeout))
	}

	if enableLog {
		db.AddQueryHook(dbLogger{})
	}

	return db, nil
}

// ConnectOptions takes a database config and returns go-pg specific connect options
func ConnectOptions(c *config.Database) *pg.Options {
	var addr = func() string {
		host := env("DB_HOST", c.Host)
		port := env("DB_PORT", c.Port)
		if c.Network == "unix" {
			if host[0:1] == "/" {
				// unix domain socket
				return host + "/.s.PGSQL." + port
			}
		}
		return host + ":" + port
	}

	return &pg.Options{
		Network:   env("DB_NETWORK", c.Network),   // "tcp" or "unix" (for unix domain sockets)
		Addr:      env("DB_ADDR", addr()),         // host:port or unix socket (e.g. /var/run/postgresql/.s.PGSQL.5432)
		Database:  env("DB_DATABASE", c.Database), // database name
		User:      env("DB_USER", c.User),         // database role
		Password:  env("DB_PASSWORD", c.Password), // plaintext password
		TLSConfig: getTLSConfig(c.SSLMode),        // "disable", "allow", "verify-ca"
	}
}

func getTLSConfig(mode string) *tls.Config {
	pgSSLMode := env("DB_SSLMODE", mode)
	if pgSSLMode == "disable" {
		return nil
	}
	// missing or any other option
	return &tls.Config{
		InsecureSkipVerify: true,
	}
}

// settings retrieves our database-specific "settings"
func (db *DB) settings() error {
	db.Settings = &sandpiper.Setting{ID: true}

	err := db.Select(db.Settings)
	if err == pg.ErrNoRows {
		return errors.New("missing db settings: database not initialized")
	}
	return err
}

func env(key, defValue string) string {
	envValue := os.Getenv(key)
	if envValue != "" {
		return envValue
	}
	return defValue
}
