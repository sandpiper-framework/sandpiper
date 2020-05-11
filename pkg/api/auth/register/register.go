// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

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
	svc := auth.Initialize(db, tgen, sec, rbac.New(db.ServerRole))
	ls := al.ServiceLogger(svc, log)
	at.NewHTTP(ls, srv, mwFunc)
}
