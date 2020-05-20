// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package command implements `sandpiper` commands (add, pull, list, ...)
package command

// sandpiper secrets

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	args "github.com/urfave/cli/v2"
)

// Secrets generates new secrets for the server config.yaml file
func Secrets(c *args.Context) error {
	fmt.Println("GENERATE RANDOM SERVER 'SECRETS'")
	api, err := APISecret()
	if err != nil {
		return err
	}
	jwt, err := JWTSecret()
	if err != nil {
		return err
	}

	fmt.Printf("\n# ENVIRONMENT VARIABLES (remove double quotes from Docker .env files):\n\n")
	fmt.Printf("export APIKEY_SECRET=\"%s\"\nexport JWT_SECRET=\"%s\"\n", api, jwt)

	fmt.Printf("\n# CONFIG FILE (YAML) ENTRIES:\n")
	fmt.Printf("\napi_key_secret: %s\nsecret: %s\n\n", api, jwt)

	return nil
}

// APISecret returns a random string suitable for an AES256 encryption key
func APISecret() (string, error) {
	// Base64 Encoded AES-256 key (44 chars)
	// similar to node -e "console.log(require('crypto').randomBytes(32).toString('base64'));"
	return SecretKey(32)
}

// JWTSecret returns a random string suitable for a JWT secret key
func JWTSecret() (string, error) {
	// similar to node -e "console.log(require('crypto').randomBytes(64).toString('base64'));"
	return SecretKey(64)
}

// SecretKey generates a random key suitable for secrets
func SecretKey(n int) (string, error) {
	s, err := randomBytes(n)
	if err != nil {
		return "", err
	}
	return toBase64(s), nil
}

func randomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}

func toBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}
