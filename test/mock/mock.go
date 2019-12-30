// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package mock

import (
	"net/http/httptest"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// TestHash returns an sha-1 hash
func TestHash(n int) string {
	s := strconv.Itoa(n)
	return s[0:1] + "a804c61e1a70ab37b912792ee846de7378c4a36"
}

// TestUUID returns a valid test uuid starting with the supplied number (1-9)
func TestUUID(n int) uuid.UUID {
	s := strconv.Itoa(n)
	return uuid.MustParse(s[0:1] + "0000000-0000-0000-0000-000000000000")
}

// TestTime is used for testing time fields
func TestTime(year int) time.Time {
	return time.Date(year, time.May, 19, 1, 2, 3, 4, time.UTC)
}

// TestTimePtr is used for testing pointer time fields
func TestTimePtr(year int) *time.Time {
	t := time.Date(year, time.May, 19, 1, 2, 3, 4, time.UTC)
	return &t
}

// HeaderValid is used for jwt testing
func HeaderValid() string {
	return "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjIjoiOWFkMDcyMzQtMjc0Mi00MGViLTllZjEtODAwYzJmMTE2NGNlIiwiZSI6ImpvaG5kb2VAbWFpbC5jb20iLCJleHAiOjE1Nzc2MzQ4NzksImlkIjoxLCJyIjoxMDAsInUiOiJhZG1pbiJ9.1i-Jwoent1Oyx0IAx4Ass7lBAjgf3O3RvxihwKCu2g4"
}

// HeaderInvalid is used for jwt testing
func HeaderInvalid() string {
	return "Bearer eyJhbGciOiJIUzM4NCIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwidSI6ImpvaG5kb2UiLCJlIjoiam9obmRvZUBtYWlsLmNvbSIsInIiOjEsImMiOjEsImwiOjEsImV4cCI6NDEwOTMyMDg5NCwiaWF0IjoxNTE2MjM5MDIyfQ.7uPfVeZBkkyhICZSEINZfPo7ZsaY0NNeg0ebEGHuAvNjFvoKNn8dWYTKaZrqE1X4"
}

// EchoCtxWithKeys returns new Echo context with keys
func EchoCtxWithKeys(keys []string, values ...interface{}) echo.Context {
	e := echo.New()
	w := httptest.NewRecorder()
	c := e.NewContext(nil, w)
	for i, k := range keys {
		c.Set(k, values[i])
	}
	return c
}
