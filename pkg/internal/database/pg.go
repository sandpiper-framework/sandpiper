// Copyright Auto Care Association. All rights reserved.
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
	// DB adapter
	_ "github.com/lib/pq"
)

type dbLogger struct{}

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
func New(psn string, timeout int, enableLog bool) (*pg.DB, error) {
	u, err := pg.ParseURL(psn)
	if err != nil {
		return nil, err
	}

	db := pg.Connect(u)

	_, err = db.Exec("SELECT 1") // test connectivity
	if err != nil {
		return nil, err
	}

	if timeout > 0 {
		db.WithTimeout(time.Second * time.Duration(timeout))
	}

	if enableLog {
		db.AddQueryHook(dbLogger{})
	}

	return db, nil
}
