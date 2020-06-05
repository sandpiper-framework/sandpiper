// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

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
