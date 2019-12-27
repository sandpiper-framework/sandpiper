// Package jwt handles json-web-token (RFC7519) middleware to represent
// claims securely between two parties.
package jwt

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/internal/model"
)

// Custom errors
var (
	// ErrInvalidToken indicates a missing or invalid token was supplied
	ErrInvalidToken = errors.New("token is missing or invalid")
)

// New generates new JWT service necessary for auth middleware
func New(secret, algo string, d int) *Service {
	signingMethod := jwt.GetSigningMethod(algo)
	if signingMethod == nil {
		panic("invalid jwt signing method")
	}
	return &Service{
		key:      []byte(secret),
		algo:     signingMethod,
		duration: time.Duration(d) * time.Minute,
	}
}

// Service provides a Json-Web-Token authentication implementation
type Service struct {
	// Secret key used for signing.
	key []byte

	// Duration for which the jwt token is valid.
	duration time.Duration

	// JWT signing algorithm
	algo jwt.SigningMethod
}

// MWFunc makes JWT implement the Middleware interface.
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
			companyID := int(claims["c"].(float64))
			username := claims["u"].(string)
			email := claims["e"].(string)
			role := sandpiper.AccessRole(claims["r"].(float64))

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

// GenerateToken generates new JWT token and populates it with user data
func (j *Service) GenerateToken(u *sandpiper.User) (string, string, error) {
	expire := time.Now().Add(j.duration)

	token := jwt.NewWithClaims(j.algo, jwt.MapClaims{
		"id":  u.ID,
		"u":   u.Username,
		"e":   u.Email,
		"r":   u.Role.AccessLevel,
		"c":   u.CompanyID,
		"exp": expire.Unix(),
	})

	tokenString, err := token.SignedString(j.key)

	return tokenString, expire.Format(time.RFC3339), err
}
