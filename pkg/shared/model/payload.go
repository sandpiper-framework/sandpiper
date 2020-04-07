// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package sandpiper

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
)

/* Utility routines to support our encoding types
 * Usage:
 *   import "autocare.org/sandpiper/pkg/shared/payload"
 *
 *   data := bytes.NewReader([]byte("payload data to store"))
 *   payloadData, err := sandpiper.Encode(data, "b64")
 *
 *   rawBytes, err := payloadData.Decode()
 */

// PayloadData is the data type for encoded payload data.
type PayloadData string

// PayloadNil is the zero value for the PayloadData type
const PayloadNil = ""

// Encode payload data for transmission and storage
func Encode(b io.Reader, enc string) (PayloadData, error) {
	if enc == "raw" {
		// no conversion, just return original from reader
		buf, err := ioutil.ReadAll(b)
		if err != nil {
			return PayloadNil, err
		}
		return PayloadData(buf), nil
	}

	if enc == "b64" {
		// convert to base64 using a new byte slice
		buf, err := ioutil.ReadAll(b)
		if err != nil {
			return PayloadNil, err
		}
		b64 := base64.StdEncoding.EncodeToString(buf)
		return PayloadData(b64), nil
	}

	if enc == "z64" {
		// compress into zipped
		var zipped bytes.Buffer
		gz, _ := gzip.NewWriterLevel(&zipped, gzip.BestCompression)
		if _, err := io.Copy(gz, b); err != nil {
			return PayloadNil, err
		}
		if err := gz.Flush(); err != nil {
			return PayloadNil, err
		}
		if err := gz.Close(); err != nil {
			return PayloadNil, err
		}
		// base64 using a new byte slice
		b64 := make([]byte, base64.StdEncoding.EncodedLen(zipped.Len()))
		base64.StdEncoding.Encode(b64, zipped.Bytes())
		return PayloadData(b64), nil
	}

	return PayloadNil, fmt.Errorf("unknown encoding \"%s\"", enc)
}

// Decode method converts encoded payload to human-readable
func (p PayloadData) Decode(enc string) (PayloadData, error) {
	if enc == "raw" {
		return p, nil
	}
	if enc == "b64" {
		b := make([]byte, base64.StdEncoding.DecodedLen(len(p)))
		_, err := base64.StdEncoding.Decode(b, []byte(p))
		return PayloadData(b), err
	}
	if enc == "z64" {
		// convert base64 to compressed binary
		z := make([]byte, base64.StdEncoding.DecodedLen(len(p)))
		_, err := base64.StdEncoding.Decode(z, []byte(p))
		if err != nil {
			return PayloadNil, err
		}
		// decompress zipped binary to original text
		reader, err := gzip.NewReader(bytes.NewReader(z))
		if err != nil {
			return PayloadNil, err
		}
		defer reader.Close()
		b, err := ioutil.ReadAll(reader)
		return PayloadData(b), err
	}
	return PayloadNil, fmt.Errorf("unknown encoding \"%s\"", enc)
}
