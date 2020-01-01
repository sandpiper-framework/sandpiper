// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package password

import (
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/internal/model"

	"autocare.org/sandpiper/pkg/api/password"
	pl "autocare.org/sandpiper/pkg/api/password/logging"
	pt "autocare.org/sandpiper/pkg/api/password/transport"
)

// Register ties the company service to its logger and transport mechanisms
func Register(db *pg.DB, rbac password.RBAC, sec password.Securer, log sandpiper.Logger, v1 *echo.Group) {
	svc := password.Initialize(db, rbac, sec)
  ls := pl.ServiceLogger(svc, log)
	pt.NewHTTP(ls, v1)
}