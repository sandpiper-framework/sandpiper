// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package payload handles payload from/to file functions
package payload

import (
	"os"

	"autocare.org/sandpiper/pkg/shared/model"
)

// FromFile encodes a filesystem file to binary data for storing in the database
func FromFile(fileName string) (sandpiper.PayloadData, error) {
	// get a reader for the file to add
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// encode file contents for grain's payload
	payload, err := sandpiper.Encode(file)
	if err != nil {
		return nil, err
	}
	return payload, nil
}
