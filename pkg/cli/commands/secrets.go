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
	api, err := apiSecret()
	if err != nil {
		return err
	}
	jwt, err := jwtSecret()
	if err != nil {
		return err
	}

	fmt.Printf("\n# ENVIRONMENT VARIABLES (remove double quotes from Docker .env files):\n\n")
	fmt.Printf("export APIKEY_SECRET=\"%s\"\nexport JWT_SECRET=\"%s\"\n", api, jwt)

	fmt.Printf("\n# CONFIG FILE (YAML) ENTRIES:\n")
	fmt.Printf("\napi_key_secret: %s\nsecret: %s\n\n", api, jwt)

	return nil
}

func apiSecret() (string, error) {
	// Base64 Encoded AES-256 key (44 chars)
	// node -e "console.log(require('crypto').randomBytes(32).toString('base64'));"
	return secretKey(32)
}

func jwtSecret() (string, error) {
	// node -e "console.log(require('crypto').randomBytes(64).toString('base64'));"
	return secretKey(64)
}

func secretKey(n int) (string, error) {
	s, err := randomBytes(n)
	if err != nil {
		return "", err
	}
	return toBase64(s), nil
}

func randomBytes(n int) (string, error) {
	b := make([]byte, n/2)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}

func toBase64(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}
