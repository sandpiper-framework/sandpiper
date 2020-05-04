// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package mock

import (
	"sandpiper/pkg/shared/model"
)

// JWT mock
type JWT struct {
	GenerateTokenFn func(*sandpiper.User) (string, string, error)
}

// GenerateToken mock
func (j *JWT) GenerateToken(u *sandpiper.User) (string, string, error) {
	return j.GenerateTokenFn(u)
}
