// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

// Package secure handles password scoring, encrypting and token generation.
package secure_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sandpiper-framework/sandpiper/pkg/shared/secure"
)

func TestEncrypt(t *testing.T) {
	cases := []struct {
		name    string
		data    []byte
		key     string
		result  []byte
		wantErr error
	}{
		{
			name:    "Back and Forth",
			data:    []byte("sandpiper rocks!"),
			key:     "u7WJ3kpqyvAkKb7HIfYJoSok2DoqTa9YhaCUhUujqb8=",
			result:  []byte("sandpiper rocks!"),
			wantErr: nil,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			e, err := secure.Encrypt(tt.data, tt.key)
			assert.Equal(t, err, tt.wantErr)
			d, err := secure.Decrypt(e, tt.key)
			assert.Equal(t, err, tt.wantErr)
			assert.Equal(t, tt.result, d)
		})
	}
}
