// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package mock

// todo: change to https://github.com/ory/dockertest

import (
	"database/sql"
	"testing"

	"github.com/fortytw2/dockertest"
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"

	"sandpiper/pkg/shared/database"
)

// NewPGContainer instantiates new postgresql docker container
func NewPGContainer(t *testing.T) *dockertest.Container {
	container, err := dockertest.RunContainer(
		"postgres:alpine",
		"5432",
		func(addr string) error {
			db, err := sql.Open("postgres", "postgres://postgres:postgres@"+addr+"?sslmode=disable")
			fatalErr(t, err)
			return db.Ping()
		},
	)
	fatalErr(t, err)
	return container
}

// NewDB instantiates new postgresql database connection via docker container
func NewDB(t *testing.T, con *dockertest.Container, models ...interface{}) *pg.DB {
	db, err := database.New("postgres://postgres:postgres@"+con.Addr+"/postgres?sslmode=disable", 10, false)
	fatalErr(t, err)

	for _, v := range models {
		fatalErr(t, db.CreateTable(v, &orm.CreateTableOptions{FKConstraints: true}))
	}

	return db
}

// InsertMultiple inserts multiple values into database
func InsertMultiple(db *pg.DB, models ...interface{}) error {
	for _, v := range models {
		if err := db.Insert(v); err != nil {
			return err
		}
	}
	return nil
}

func fatalErr(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
