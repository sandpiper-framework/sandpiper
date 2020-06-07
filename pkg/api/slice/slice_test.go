// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package slice_test

import (
	"github.com/sandpiper-framework/sandpiper/pkg/api/slice"
	"testing"
)

func TestCreate(t *testing.T) {

}

func TestView(t *testing.T) {

}

func TestList(t *testing.T) {

}

func TestDelete(t *testing.T) {

}

func TestUpdate(t *testing.T) {

}

func TestInitialize(t *testing.T) {
	s := slice.Initialize(nil, nil, nil)
	if s == nil {
		t.Error("Slice service not initialized")
	}
}
