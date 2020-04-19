// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package client is used to communicate with the api server
package client

// extracted with slight mods from github.com/ddliu/go-httpclient/httpclient.go

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
)

// Response is a thin wrapper of http.Response (can also be used as http.Response)
// allowing us to add new methods
type Response struct {
	*http.Response
}

// ReadAll response body into a byte slice.
func (r *Response) ReadAll() ([]byte, error) {
	var reader io.ReadCloser
	var err error

	switch r.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(r.Body)
		if err != nil {
			return nil, err
		}
	default:
		reader = r.Body
	}
	defer reader.Close()
	return ioutil.ReadAll(reader)
}

// ToString reads response body into string.
func (r *Response) ToString() (string, error) {
	bytes, err := r.ReadAll()
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
