// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package sandpiper

import (
	//"encoding/base64"
	"bytes"
	"compress/gzip"
	"database/sql/driver"
	"errors"
	"io/ioutil"

	//"github.com/go-pg/pg/v9"
	// DB adapter
	_ "github.com/lib/pq"
)

// PayloadType is a custom data type for payload data to make encoding/decoding
// the payload data automatic.
type PayloadType []byte

// Value implements driver.Valuer
// Processes data into the database
func (p PayloadType) Value() (driver.Value, error) {
	b := make([]byte, 0, len(p))
	buf := bytes.NewBuffer(b)
	w, _ := gzip.NewWriterLevel(buf, gzip.BestCompression)
	_, err := w.Write(p)
	if err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return driver.Value(buf.Bytes()), nil
}

// Scan implements sql.Scanner. Only handles string and []byte
// Processes data from the database
func (p *PayloadType) Scan(src interface{}) error {
	var b []byte

	switch src.(type) {
	case string:
		b = []byte(src.(string))
	case []byte:
		b = src.([]byte)
	default:
		return errors.New("PayloadType must be 'string' or '[]byte'")
	}

	reader, _ := gzip.NewReader(bytes.NewReader(b))
	defer reader.Close()
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	*p = PayloadType(b)
	return nil
}
