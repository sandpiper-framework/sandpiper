// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package credentials handles authorization with primary server
package credentials

import (
	"encoding/hex"
	"encoding/json"
	"errors"

	"sandpiper/pkg/shared/secure"
)

// SyncLogin is used to manage credentials for the sync_api_key
type SyncLogin struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

// New returns a SyncLogin decoded from the supplied api_key (base64) and secret
func New(APIKey, secret string) (*SyncLogin, error) {
	b, err := fromHex(APIKey)
	if err != nil {
		return nil, err
	}
	data, err := secure.Decrypt(b, secret)
	if err != nil {
		return nil, err
	}
	login := new(SyncLogin)
	if err := json.Unmarshal(data, login); err != nil {
		return nil, errors.New("improper format of api-key")
	}
	return login, nil
}

// APIKey returns an encrypted sync_api_key (base64) from sync login credentials
func (s *SyncLogin) APIKey(secret string) ([]byte, error) {
	login, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	data, err := secure.Encrypt(login, secret)
	if err != nil {
		return nil, err
	}
	return toHex(data), nil
}

func toHex(src []byte) []byte {
	buf := make([]byte, hex.EncodedLen(len(src)))
	_ = hex.Encode(buf, src)
	return buf
}

func fromHex(s string) ([]byte, error) {
	return hex.DecodeString(s)
}
