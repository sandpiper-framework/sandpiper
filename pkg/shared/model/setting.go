// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package sandpiper

// Setting represents the setting domain model
type Setting struct {
	ID    int    `json:"id" pg:",pk"`
	Key   string `json:"key"`
	Value string `json:"value"`
}
