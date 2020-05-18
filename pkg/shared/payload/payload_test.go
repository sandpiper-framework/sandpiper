// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package payload_test

import (
	"bytes"
	"fmt"
	"testing"

	"sandpiper/pkg/shared/payload"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		name    string
		src     []byte
		enc     string
		want    string
		wantErr bool
	}{
		{
			name:    "Good Conversion",
			src:     []byte("sandpiper rocks!"),
			enc:     "z64",
			want:    "H4sIAAAAAAAC/ypOzEspyCxILVIoyk/OLlYEAAAA//8BAAD//451mN4QAAAA",
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data := bytes.NewReader(test.src)
			got, err := payload.Encode(data, test.enc)
			if (err != nil) != test.wantErr {
				t.Errorf("error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if string(got) != test.want {
				t.Errorf("got = %s\n, want %s\n", got, test.want)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	tests := []struct {
		name    string
		data    payload.PayloadData
		enc     string
		want    string
		wantErr bool
	}{
		{
			name:    "Good Conversion",
			data:    "H4sIAAAAAAAC/ypOzEspyCxILVIoyk/OLlYEAAAA//8BAAD//451mN4QAAAA",
			enc:     "z64",
			want:    "sandpiper rocks!",
			wantErr: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.data.Decode(test.enc)
			if (err != nil) != test.wantErr {
				t.Errorf("error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if string(got) != test.want {
				t.Errorf("got = %s\n, want %s\n", got, test.want)
			}
		})
	}
}

func TestPayloadData(t *testing.T) {
	src := []byte("sandpiper rocks!")

	t.Run("Back & Forth", func(t *testing.T) {
		data := bytes.NewReader(src)
		in, err := payload.Encode(data, "b64")
		if err != nil {
			t.Errorf("error = %v", err)
			return
		}
		out, err := in.Decode("b64")
		if err != nil {
			t.Errorf("error = %v", err)
			return
		}
		if string(src) != out {
			t.Errorf("got \"%s\" want \"%s\"\n", src, out)
			b := []byte(out)
			fmt.Printf("out: %v", b)
			return
		}
	})
}

/*
func main() {
	s := []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.")

	for i := 10; i <= 440; i = i + 10 {
		s2, _ := encode(s[:i])
		fmt.Printf("input: %d, output %d\n", i, len(s2))
	}
}
*/
