// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package database_test

import (
	"database/sql"
	"testing"

	"github.com/fortytw2/dockertest"
	"github.com/stretchr/testify/assert"

	"github.com/sandpiper-framework/sandpiper/pkg/shared/database"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/model"
)

func TestNew(t *testing.T) {
	container, err := dockertest.RunContainer("postgres:alpine", "5432", func(addr string) error {
		db, err := sql.Open("postgres", "postgres://postgres:postgres@"+addr+"?sslmode=disable")
		if err != nil {
			return err
		}

		return db.Ping()
	})
	defer container.Shutdown()
	if err != nil {
		t.Fatalf("could not start postgres, %s", err)
	}

	_, err = database.New("PSN", 1, false)
	if err == nil {
		t.Error("Expected error")
	}

	_, err = database.New("postgres://postgres:postgres@localhost:1234/postgres?sslmode=disable", 0, false)
	if err == nil {
		t.Error("Expected error")
	}

	dbLogTest, err := database.New("postgres://postgres:postgres@"+container.Addr+"/postgres?sslmode=disable", 0, true)
	if err != nil {
		t.Fatalf("Error establishing connection %v", err)
	}
	dbLogTest.Close()

	db, err := database.New("postgres://postgres:postgres@"+container.Addr+"/postgres?sslmode=disable", 1, true)
	if err != nil {
		t.Fatalf("Error establishing connection %v", err)
	}

	var user sandpiper.User
	db.Select(&user)

	assert.NotNil(t, db)

	db.Close()

}
