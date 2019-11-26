package postgres_test

import (
	"database/sql"
	"testing"

	"github.com/fortytw2/dockertest"
	"github.com/stretchr/testify/assert"

	"autocare.org/sandpiper/pkg/model"
	"autocare.org/sandpiper/pkg/postgres"
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

	_, err = postgres.New("PSN", 1, false)
	if err == nil {
		t.Error("Expected error")
	}

	_, err = postgres.New("postgres://postgres:postgres@localhost:1234/postgres?sslmode=disable", 0, false)
	if err == nil {
		t.Error("Expected error")
	}

	dbLogTest, err := postgres.New("postgres://postgres:postgres@"+container.Addr+"/postgres?sslmode=disable", 0, true)
	if err != nil {
		t.Fatalf("Error establishing connection %v", err)
	}
	dbLogTest.Close()

	db, err := postgres.New("postgres://postgres:postgres@"+container.Addr+"/postgres?sslmode=disable", 1, true)
	if err != nil {
		t.Fatalf("Error establishing connection %v", err)
	}

	var user sandpiper.User
	db.Select(&user)

	assert.NotNil(t, db)

	db.Close()

}
