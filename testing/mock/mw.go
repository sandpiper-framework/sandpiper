package mock

import (
	"autocare.org/sandpiper/pkg/model"
)

// JWT mock
type JWT struct {
	GenerateTokenFn func(*sandpiper.User) (string, string, error)
}

// GenerateToken mock
func (j *JWT) GenerateToken(u *sandpiper.User) (string, string, error) {
	return j.GenerateTokenFn(u)
}
