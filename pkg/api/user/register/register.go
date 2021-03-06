// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package user

import (
	"github.com/labstack/echo/v4"

	"github.com/sandpiper-framework/sandpiper/pkg/api/user"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/database"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/model"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/rbac"

	ul "github.com/sandpiper-framework/sandpiper/pkg/api/user/logging"
	ut "github.com/sandpiper-framework/sandpiper/pkg/api/user/transport"
)

// Register ties the user service to its logger and transport mechanisms
func Register(db *database.DB, sec user.Securer, log sandpiper.Logger, v1 *echo.Group) {
	rba := rbac.New(db.Settings.ServerRole)
	rba.ServerID = db.Settings.ServerID
	svc := user.Initialize(db, rba, sec)
	ls := ul.ServiceLogger(svc, log)
	ut.NewHTTP(ls, v1)
}
