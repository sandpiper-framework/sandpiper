// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

// Package secure handles security and login routines
package secure

import (
	"encoding/hex"
	"encoding/json"
	"errors"
)

// Credentials is used to manage credentials for user and sync logins
type Credentials struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	SyncAPIKey string `json:"sync-api-key"`
}

// NewCredentials returns Credentials decoded from the supplied api_key (base64) and secret
func NewCredentials(APIKey, secret string) (*Credentials, error) {
	b, err := fromHex(APIKey)
	if err != nil {
		return nil, err
	}
	data, err := Decrypt(b, secret)
	if err != nil {
		return nil, err
	}
	creds := new(Credentials)
	if err := json.Unmarshal(data, creds); err != nil {
		return nil, errors.New("improper format of api-key")
	}
	return creds, nil
}

// APIKey returns an encrypted sync_api_key (base64) from sync login credentials
func (s *Credentials) APIKey(secret string) ([]byte, error) {
	login, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	data, err := Encrypt(login, secret)
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
