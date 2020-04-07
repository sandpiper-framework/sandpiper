// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package payload handles payload from/to file functions
package payload

import (
	"os"

	"autocare.org/sandpiper/pkg/shared/model"
)

// FromFile encodes a filesystem file for storing in the database
func FromFile(fileName string) (sandpiper.PayloadData, error) {
	// get a reader for the file to add
	file, err := os.Open(fileName)
	if err != nil {
		return sandpiper.PayloadNil, err
	}
	defer file.Close()

	// encode file contents for grain's payload
	payload, err := sandpiper.Encode(file, "z64")
	if err != nil {
		return sandpiper.PayloadNil, err
	}
	return payload, nil
}
