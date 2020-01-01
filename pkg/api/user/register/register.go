// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package user

import (
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/internal/model"

	"autocare.org/sandpiper/pkg/api/user"
	ul "autocare.org/sandpiper/pkg/api/user/logging"
	ut "autocare.org/sandpiper/pkg/api/user/transport"
)

// Register ties the company service to its logger and transport mechanisms
func Register(db *pg.DB, rbac user.RBAC, sec user.Securer, log sandpiper.Logger, v1 *echo.Group) {
	svc := user.Initialize(db, rbac, sec)
  ls := ul.ServiceLogger(svc, log)
	ut.NewHTTP(ls, v1)
}