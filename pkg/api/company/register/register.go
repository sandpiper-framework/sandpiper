// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package company

import (
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/api/company"
	"autocare.org/sandpiper/pkg/shared/model"
	"autocare.org/sandpiper/pkg/shared/rbac"

	cl "autocare.org/sandpiper/pkg/api/company/logging"
	ct "autocare.org/sandpiper/pkg/api/company/transport"
)

// Register ties the company service to its logger and transport mechanisms
func Register(db *pg.DB, sec company.Securer, log sandpiper.Logger, v1 *echo.Group) {
	rba := rbac.New()
	rba.ScopingField = "id" // company service so scoping is simply "id"
	svc := company.Initialize(db, rba, sec)
	ls := cl.ServiceLogger(svc, log)
	ct.NewHTTP(ls, v1)
}
