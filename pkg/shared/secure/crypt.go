// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

// Package secure handles password scoring, encrypting and token generation.
package secure

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

// The following encrypt/decrypt functions are based partially on:
// https://www.thepolyglotdeveloper.com/2018/02/encrypt-decrypt-data-golang-application-crypto-packages/

// Encrypt creates ciphered binary data from supplied plain-text data using AES-256
func Encrypt(data []byte, keyB64 string) ([]byte, error) {

	key, err := binaryKey(keyB64)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	cipherText := gcm.Seal(nonce, nonce, data, nil)

	return cipherText, nil
}

// Decrypt takes data encrypted with Encrypt function and returns the original decrypted data.
// Must use the original base64 aes-key
func Decrypt(data []byte, keyB64 string) ([]byte, error) {

	key, err := binaryKey(keyB64)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	n := gcm.NonceSize()
	nonce, cipherText := data[:n], data[n:]

	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}

func binaryKey(keyB64 string) ([]byte, error) {
	// decode key from base64 to binary (should be 256-bits)
	key, err := base64.StdEncoding.DecodeString(keyB64)
	if err != nil {
		return nil, err
	}
	return key, nil
}
