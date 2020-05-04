// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package slice_test

import (
	"sandpiper/pkg/api/slice"
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
