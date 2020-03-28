// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package secure handles password scoring, encrypting and token generation.
package secure

import (
	"crypto/sha1"
	"fmt"
	"hash"
	"strconv"
	"time"

	"github.com/nbutton23/zxcvbn-go"
	"golang.org/x/crypto/bcrypt"
)

// Service holds security related info.
type Service struct {
	minScore int // 0,1,2,3,4
	h        hash.Hash
}

// New initializes security service.
func New(minPasswordScore int) *Service {
	return &Service{minScore: minPasswordScore, h: sha1.New()}
}

// Password checks whether password is secure enough using zxcvbn library.
func (s *Service) Password(pass string, inputs ...string) bool {
	pw := zxcvbn.PasswordStrength(pass, inputs)
	return pw.Score >= s.minScore
}

// Hash encrypts the password using bcrypt.
func (*Service) Hash(password string) string {
	hashedPW, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPW)
}

// HashMatchesPassword matches hash with password. Returns true if hash and password match.
func (*Service) HashMatchesPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// Token generates new unique token.
func (s *Service) Token(str string) string {
	s.h.Reset()
	_, _ = fmt.Fprintf(s.h, "%s%s", str, strconv.Itoa(time.Now().Nanosecond()))
	return fmt.Sprintf("%x", s.h.Sum(nil))
}
