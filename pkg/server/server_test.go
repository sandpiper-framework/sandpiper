package server_test

import (
	"testing"

	"autocare.org/sandpiper/pkg/server"
)

// Improve tests
func TestNew(t *testing.T) {
	e := server.New()
	if e == nil {
		t.Errorf("Server should not be nil")
	}
}
