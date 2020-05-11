// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package company

import (
	"github.com/labstack/echo/v4"

	"sandpiper/pkg/api/company"
	"sandpiper/pkg/shared/model"
	"sandpiper/pkg/shared/rbac"

	cl "sandpiper/pkg/api/company/logging"
	ct "sandpiper/pkg/api/company/transport"
	"sandpiper/pkg/shared/database"
)

// Register ties the company service to its logger and transport mechanisms
func Register(db *database.DB, sec company.Securer, log sandpiper.Logger, v1 *echo.Group) {
	rba := rbac.New(db.ServerRole)
	rba.ScopingField = "id" // company service so scoping is simply "id"
	svc := company.Initialize(db, rba, sec)
	ls := cl.ServiceLogger(svc, log)
	ct.NewHTTP(ls, v1)
}
