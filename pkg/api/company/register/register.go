// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package company

import (
	"github.com/labstack/echo/v4"

	"github.com/sandpiper-framework/sandpiper/pkg/api/company"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/model"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/rbac"

	cl "github.com/sandpiper-framework/sandpiper/pkg/api/company/logging"
	ct "github.com/sandpiper-framework/sandpiper/pkg/api/company/transport"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/database"
)

// Register ties the company service to its logger and transport mechanisms
func Register(db *database.DB, sec company.Securer, log sandpiper.Logger, v1 *echo.Group) {
	rba := rbac.New(db.Settings.ServerRole)
	rba.ScopingField = "id" // company service so scoping is simply "id"
	svc := company.Initialize(db, rba, sec)
	ls := cl.ServiceLogger(svc, log)
	ct.NewHTTP(ls, v1)
}
