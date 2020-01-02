// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package company

import (
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/internal/model"

	"autocare.org/sandpiper/pkg/api/slice"
	sl "autocare.org/sandpiper/pkg/api/slice/logging"
	st "autocare.org/sandpiper/pkg/api/slice/transport"
)

// Register ties the company service to its logger and transport mechanisms
func Register(db *pg.DB, rbac slice.RBAC, sec slice.Securer, log sandpiper.Logger, v1 *echo.Group) {
	svc := slice.Initialize(db, rbac, sec)
	ls := sl.ServiceLogger(svc, log)
	st.NewHTTP(ls, v1)
}
