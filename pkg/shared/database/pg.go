// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package database creates a pooled connection to the database. We use a
// lightweight ORM (with deep support for postgresql). This ORM only supports
// postgresql. We might consider switching if require support for other dbms.
package database

import (
	"context"
	"log"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/google/uuid"

	// DB adapter
	_ "github.com/lib/pq"

	"sandpiper/pkg/shared/model"
)

type dbLogger struct{}

// DB is a wrapper around pg.DB so we can add functionality
type DB struct {
	*pg.DB
	Settings map[string]string
	// the following are required settings
	ServerRole string
	ServerID   uuid.UUID
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
func New(psn string, timeout int, enableLog bool) (*DB, error) {
	uri, err := pg.ParseURL(psn)
	if err != nil {
		return nil, err
	}

	// wrap db connection in our struct
	db := &DB{DB: pg.Connect(uri)}

	// test connectivity
	_, err = db.Exec("SELECT 1")
	if err != nil {
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

// settings retrieves any key/value pairs from the database "settings" table.
func (db *DB) settings() error {
	var settings []sandpiper.Setting

	err := db.Model(&settings).Select()
	if err != nil && err != pg.ErrNoRows {
		return err
	}

	db.Settings = make(map[string]string)
	for _, setting := range settings {
		db.Settings[setting.Key] = setting.Value
	}

	return nil
}
