// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package secure_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sandpiper-framework/sandpiper/pkg/shared/secure"
)

// We can't test the results one-way because the AES algorithm is not deterministic. (A new
// "Initialization Vector" is used each time)

func TestCreds(t *testing.T) {
	tests := []struct {
		name   string
		login  *secure.Credentials
		secret string
	}{
		{
			name: "Back & Forth",
			login: &secure.Credentials{
				Username: "mickeymouse",
				Password: "minnie",
			},
			secret: "u7WJ3kpqyvAkKb7HIfYJoSok2DoqTa9YhaCUhUujqb8=",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			key, err := test.login.APIKey(test.secret)
			assert.Equal(t, err != nil, false)
			c, err := secure.NewCredentials(string(key), test.secret)
			assert.Equal(t, err != nil, false)
			assert.Equal(t, c != nil, false)
			if c != nil {
				assert.Equal(t, test.login.Username, c.Username)
				assert.Equal(t, test.login.Password, c.Password)
			}
		})
	}
}
