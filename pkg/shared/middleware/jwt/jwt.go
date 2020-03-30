// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package jwt handles json-web-token (RFC7519) middleware to represent
// claims securely between two parties.
package jwt

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/shared/model"
)

// Custom errors
var (
	// ErrInvalidToken indicates a missing or invalid token was supplied
	ErrInvalidToken = errors.New("token is missing or invalid")
)

// New generates new JWT service necessary for auth middleware
func New(secret, algo string, ttlMinutes int, minLen int) (*Service, error) {
	var minSecretLen = 32 // based on guidelines using the HS256 algorithm

	if minLen > 0 { // check for an override of the minimum value
		minSecretLen = minLen
	}
	if len(secret) < minSecretLen {
		return nil, fmt.Errorf("jwt secret length is %v, which is less than required %v", len(secret), minSecretLen)
	}
	signingMethod := jwt.GetSigningMethod(algo)
	if signingMethod == nil {
		return nil, fmt.Errorf("invalid jwt signing method: %s", algo)
	}
	return &Service{
		key:  []byte(secret),
		algo: signingMethod,
		ttl:  time.Duration(ttlMinutes) * time.Minute,
	}, nil
}

// Service provides a Json-Web-Token authentication implementation
type Service struct {
	// Secret key used for signing.
	key []byte

	// Duration for which the jwt token is valid.
	ttl time.Duration

	// JWT signing algorithm
	algo jwt.SigningMethod
}

// MWFunc makes JWT implement the Middleware interface.
// Add jwt claims to context on every request.
func (j *Service) MWFunc() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := j.ParseToken(c)
			if err != nil || !token.Valid {
				return c.NoContent(http.StatusUnauthorized)
			}

			// get the JWT claims
			claims := token.Claims.(jwt.MapClaims)
			id := int(claims["id"].(float64))
			companyID, _ := uuid.Parse(claims["c"].(string))
			username := claims["u"].(string)
			email := claims["e"].(string)
			role := sandpiper.AccessLevel(claims["r"].(float64))

			// add claims to context
			c.Set("id", id)
			c.Set("company_id", companyID)
			c.Set("username", username)
			c.Set("email", email)
			c.Set("role", role)

			return next(c)
		}
	}
}

// ParseToken parses token from Authorization header
func (j *Service) ParseToken(c echo.Context) (*jwt.Token, error) {

	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return nil, ErrInvalidToken
	}
	parts := strings.SplitN(token, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return nil, ErrInvalidToken
	}

	return jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
		if j.algo != token.Method {
			return nil, ErrInvalidToken
		}
		return j.key, nil
	})
}

// GenerateToken generates new JWT token and populates it with user data for rbac
func (j *Service) GenerateToken(u *sandpiper.User) (string, string, error) {
	expire := time.Now().Add(j.ttl)

	token := jwt.NewWithClaims(j.algo, jwt.MapClaims{
		"id":  u.ID,
		"u":   u.Username,
		"e":   u.Email,
		"r":   u.Role,
		"c":   u.CompanyID,
		"exp": expire.Unix(),
	})

	tokenString, err := token.SignedString(j.key)

	return tokenString, expire.Format(time.RFC3339), err
}
