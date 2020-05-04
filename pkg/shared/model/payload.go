// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package sandpiper

import (
	"bytes"
	"compress/gzip"
	"encoding/ascii85"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
)

/* Utility routines to support our encoding types
 * Usage:
 *   import "sandpiper/pkg/shared/payload"
 *
 *   data := bytes.NewReader([]byte("payload data to store"))
 *   payloadData, err := sandpiper.Encode(data, "b64")
 *
 *   rawBytes, err := payloadData.Decode()
 */

// todo: change from gzip library to https://github.com/klauspost/compress

// PayloadData is the data type for encoded payload data.
type PayloadData string

// PayloadNil is the zero value for the PayloadData type
const PayloadNil = ""

// Encode payload data for transmission and storage
func Encode(b io.Reader, enc string) (PayloadData, error) {
	switch enc {
	case "raw":
		// no conversion, just return original from reader
		buf, err := ioutil.ReadAll(b)
		if err != nil {
			return PayloadNil, err
		}
		return PayloadData(buf), nil
	case "a85":
		// convert to ascii85 (1.25 size)
		buf, err := ioutil.ReadAll(b)
		if err != nil {
			return PayloadNil, err
		}
		return PayloadData(toAscii85(buf)), nil
	case "b64":
		// convert to base64 (1.33 size)
		buf, err := ioutil.ReadAll(b)
		if err != nil {
			return PayloadNil, err
		}
		return PayloadData(toBase64(buf)), nil
	case "z64":
		// compress and encode base64
		gz, err := toZip(b)
		if err != nil {
			return PayloadNil, err
		}
		return PayloadData(toBase64(gz)), nil
	case "z85":
		// compress and encode ascii85
		gz, err := toZip(b)
		if err != nil {
			return PayloadNil, err
		}
		return PayloadData(toAscii85(gz)), nil
	default:
		return PayloadNil, fmt.Errorf("unknown encoding \"%s\"", enc)
	}
}

func toZip(b io.Reader) ([]byte, error) {
	var zipped bytes.Buffer
	gz, _ := gzip.NewWriterLevel(&zipped, gzip.BestCompression)
	if _, err := io.Copy(gz, b); err != nil {
		return nil, err
	}
	if err := gz.Flush(); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return zipped.Bytes(), nil
}

func toBase64(buf []byte) []byte {
	b64 := make([]byte, base64.StdEncoding.EncodedLen(len(buf)))
	base64.StdEncoding.Encode(b64, buf)
	return b64
}

func toAscii85(buf []byte) []byte {
	a85 := make([]byte, ascii85.MaxEncodedLen(len(buf)))
	ascii85.Encode(a85, buf)
	return a85
}

// Decode method converts encoded payload to human-readable
func (p PayloadData) Decode(enc string) (PayloadData, error) {
	// todo: change to use io.writer? (want to avoid copying payload in memory)

	switch enc {
	case "raw":
		return p, nil
	case "a85":
		buf, err := fromAscii85([]byte(p))
		return PayloadData(buf), err
	case "b64":
		buf, err := fromBase64([]byte(p))
		return PayloadData(buf), err
	case "z85":
		// convert ascii85 to compressed binary to original
		gz, err := fromAscii85([]byte(p))
		if err != nil {
			return PayloadNil, err
		}
		buf, err := fromGzip(gz)
		return PayloadData(buf), err
	case "z64":
		// convert base64 to compressed binary to original
		gz, err := fromBase64([]byte(p))
		if err != nil {
			return PayloadNil, err
		}
		buf, err := fromGzip(gz)
		return PayloadData(buf), err
	}
	return PayloadNil, fmt.Errorf("unknown encoding \"%s\"", enc)
}

func fromAscii85(a85 []byte) ([]byte, error) {
	maxLen := int64(float64(len(a85)) * .85) //80% efficient
	buf := make([]byte, maxLen)
	_, _, err := ascii85.Decode(buf, a85, true)
	return buf, err
}

func fromBase64(b64 []byte) ([]byte, error) {
	buf := make([]byte, base64.StdEncoding.DecodedLen(len(b64)))
	_, err := base64.StdEncoding.Decode(buf, b64)
	return buf, err
}

func fromGzip(gz []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(gz))
	if err != nil {
		return nil, err
	}
	// This will avoid invalid header errors (default expects multiple files in the stream)
	reader.Multistream(false)

	defer reader.Close()
	return ioutil.ReadAll(reader)
}

/*
NOTE: string([]byte) makes a copy... Might need to use this unsafe method to avoid
the memory if performance is a problem! (What happens with a 200MB payload?)

func BytesToString(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{bh.Data, bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}
*/
