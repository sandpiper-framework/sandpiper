// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package sandpiper_test

import (
	"autocare.org/sandpiper/pkg/internal/model"
	"testing"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		name string
		src []byte
		want string
		wantErr bool
	}{
		{"Good Conversion", []byte("sandpiper rocks!"), "H4sIAAAAAAAC/ypOzEspyCxILVIoyk/OLlYEAAAA//8BAAD//451mN4QAAAA", false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := sandpiper.Encode(test.src)
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
		name string
		data sandpiper.PayloadData
		want string
		wantErr bool
	}{
		{"Good Conversion", []byte("H4sIAAAAAAAC/ypOzEspyCxILVIoyk/OLlYEAAAA//8BAAD//451mN4QAAAA"), "sandpiper rocks!", false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.data.Decode()
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
		in, err := sandpiper.Encode(src)
		if err != nil {
			t.Errorf("error = %v", err)
			return
		}
		out, err := in.Decode()
		if err != nil {
			t.Errorf("error = %v", err)
			return
		}
		if string(src) != string(out) {
			t.Errorf("got \"%s\" want \"%s\"\n", string(src), string(out))
			return
		}
	})
}