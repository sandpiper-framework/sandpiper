// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package sandpiper

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
)

/*
 * Usage:
 *   import "autocare.org/sandpiper/pkg/shared/payload"
 *
 *   rawBytes := []byte("payload data to store")
 *   payloadData, err := sandpiper.Encode(rawBytes)
 *
 *   rawBytes, err := payloadData.Decode()
 */

// PayloadData is the data type for encoded payload data.
type PayloadData []byte

// Encode payload data for transmission and storage
func Encode(b []byte) (PayloadData, error) {
	var zipped bytes.Buffer

	// zip source data into local buffer
	gz, _ := gzip.NewWriterLevel(&zipped, gzip.BestCompression)
	if _, err := gz.Write(b); err != nil {
		return nil, err
	}
	if err := gz.Flush(); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}

	// convert to base64 using a new byte slice
	b64 := make([]byte, base64.StdEncoding.EncodedLen(zipped.Len()))
	base64.StdEncoding.Encode(b64, zipped.Bytes())

	return PayloadData(b64), nil
}

// Decode method converts base64 compressed payload to byte slice
func (p PayloadData) Decode() ([]byte, error) {

	// convert base64 to compressed binary
	zipped := make([]byte, base64.StdEncoding.DecodedLen(len(p)))
	_, err := base64.StdEncoding.Decode(zipped, p)
	if err != nil {
		return nil, err
	}

	// convert compressed binary to original text
	reader, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return ioutil.ReadAll(reader)
}
