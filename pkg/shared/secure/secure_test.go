package secure_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"sandpiper/pkg/shared/secure"
)

func TestPassword(t *testing.T) {
	cases := []struct {
		name   string
		pass   string
		inputs []string
		want   bool
	}{
		{
			name: "Insecure password",
			pass: "notSec",
			want: false,
		},
		{
			name:   "Password matches input fields",
			pass:   "johndoe92",
			inputs: []string{"John", "Doe"},
			want:   false,
		},
		{
			name:   "Secure password",
			pass:   "callgophers",
			inputs: []string{"John", "Doe"},
			want:   true,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := secure.New(1)
			got := s.Password(tt.pass, tt.inputs...)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHashAndMatch(t *testing.T) {
	cases := []struct {
		name string
		pass string
		want bool
	}{
		{
			name: "Success",
			pass: "gamepad",
			want: true,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := secure.New(1)
			hash := s.Hash(tt.pass)
			assert.Equal(t, tt.want, s.HashMatchesPassword(hash, tt.pass))
		})
	}
}

func TestToken(t *testing.T) {
	s := secure.New(1)
	token := "token"
	tokenized := s.Token(token)
	assert.NotEqual(t, tokenized, token)
}
