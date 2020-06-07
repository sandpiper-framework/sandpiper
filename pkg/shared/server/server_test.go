package server_test

import (
	"testing"

	"github.com/sandpiper-framework/sandpiper/pkg/shared/server"
)

// Improve tests
func TestNew(t *testing.T) {
	e := server.New()
	if e == nil {
		t.Errorf("Server should not be nil")
	}
}
