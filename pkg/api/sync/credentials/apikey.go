// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package credentials handles authorization with primary server
package credentials

import (
	"errors"
	"strings"

	"autocare.org/sandpiper/pkg/shared/secure"
)

// SyncLogin is used to manage credentials for the sync_api_key
type SyncLogin struct {
	User     string
	Password string
}

// New returns a SyncLogin decoded from the supplied api_key and secret
func New(APIKey, secret string) (*SyncLogin, error) {
	data, err := secure.Decrypt([]byte(APIKey), secret)
	if err != nil {
		return nil, err
	}
	a := strings.Split(string(data), " ")
	if len(a) != 2 {
		return nil, errors.New("improper format of api-key")
	}
	return &SyncLogin{
		User:     a[0],
		Password: a[1],
	}, nil
}

// APIKey returns a sync_api_key from the encrypted sync login credentials
func (s *SyncLogin) APIKey(secret string) ([]byte, error) {
	b := []byte(s.User + " " + s.Password)
	data, err := secure.Encrypt(b, secret)
	if err != nil {
		return nil, err
	}
	return data, nil
}
