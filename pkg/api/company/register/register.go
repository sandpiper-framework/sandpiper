// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package company

import (
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/internal/model"

	"autocare.org/sandpiper/pkg/api/company"
	cl "autocare.org/sandpiper/pkg/api/company/logging"
	ct "autocare.org/sandpiper/pkg/api/company/transport"
)

// Register ties the company service to its logger and transport mechanisms
func Register(db *pg.DB, rbac company.RBAC, sec company.Securer, log sandpiper.Logger, v1 *echo.Group) {
	svc := company.Initialize(db, rbac, sec)
	ls := cl.ServiceLogger(svc, log)
	ct.NewHTTP(ls, v1)
}
