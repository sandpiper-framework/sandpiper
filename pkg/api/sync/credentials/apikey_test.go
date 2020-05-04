// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package credentials_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	creds "sandpiper/pkg/api/sync/credentials"
)

// We can't test the results one-way because the AES algorithm is not deterministic. (A new
// "Initialization Vector" is used each time)

func TestCreds(t *testing.T) {
	tests := []struct {
		name   string
		login  *creds.SyncLogin
		secret string
	}{
		{
			name: "Back & Forth",
			login: &creds.SyncLogin{
				User:     "mickeymouse",
				Password: "minnie",
			},
			secret: "u7WJ3kpqyvAkKb7HIfYJoSok2DoqTa9YhaCUhUujqb8=",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			key, err := test.login.APIKey(test.secret)
			assert.Equal(t, err != nil, false)
			c, err := creds.New(string(key), test.secret)
			assert.Equal(t, err != nil, false)
			assert.Equal(t, c != nil, false)
			if c != nil {
				assert.Equal(t, test.login.User, c.User)
				assert.Equal(t, test.login.Password, c.Password)
			}
		})
	}
}
