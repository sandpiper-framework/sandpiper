// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package payload handles payload from/to file functions
package payload

import (
	"os"

	"sandpiper/pkg/shared/payload"
)

// FromFile encodes a filesystem file for storing in the database
func FromFile(fileName string, enc string) (payload.PayloadData, error) {
	// get a reader for the file to add
	file, err := os.Open(fileName)
	if err != nil {
		return payload.Nil, err
	}
	defer file.Close()

	// encode file contents for grain's payload
	data, err := payload.Encode(file, enc)
	if err != nil {
		return payload.Nil, err
	}
	return data, nil
}
