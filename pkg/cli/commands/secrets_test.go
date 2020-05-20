// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package command

import (
	"encoding/base64"
	"testing"
)

// TestSecrets can't test randomness, but it does make sure the lengths are correct and the
// results are a valid base64 string
func TestSecrets(t *testing.T) {
	tests := []struct {
		name     string
		numChars int
		want     int
		wantErr  bool
	}{
		{
			name:     "AES-256 Requirements",
			numChars: 32,
			want:     44,
			wantErr:  false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := SecretKey(test.numChars)
			if (err != nil) != test.wantErr {
				t.Errorf("secretkey error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if len(got) != test.want {
				t.Errorf("got = %d\n, want %d\n", len(got), test.want)
			}
			s, err := base64.StdEncoding.DecodeString(got)
			if (err != nil) != test.wantErr {
				t.Errorf("base64 decode error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if len(s) != test.numChars {
				t.Errorf("decoded length (%d) does not match original requested chars (%d)", len(s), test.numChars)
			}
		})
	}
}
