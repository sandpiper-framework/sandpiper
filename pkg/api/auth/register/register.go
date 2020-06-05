// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package auth

import (
	"github.com/labstack/echo/v4"

	"sandpiper/pkg/api/auth"
	"sandpiper/pkg/shared/model"
	"sandpiper/pkg/shared/rbac"

	al "sandpiper/pkg/api/auth/logging"
	at "sandpiper/pkg/api/auth/transport"
	"sandpiper/pkg/shared/database"
)

// Register ties the auth service to its logger and transport mechanisms
func Register(db *database.DB, sec auth.Securer, log sandpiper.Logger, srv *echo.Echo, tgen auth.TokenGenerator, mwFunc echo.MiddlewareFunc) {
	rba := rbac.New(db.Settings.ServerRole)
	rba.ServerID = db.Settings.ServerID
	svc := auth.Initialize(db, tgen, sec, rba)
	ls := al.ServiceLogger(svc, log)
	at.NewHTTP(ls, srv, mwFunc)
}
