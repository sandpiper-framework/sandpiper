// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package auth

import (
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/api/auth"
	"autocare.org/sandpiper/pkg/shared/model"
	"autocare.org/sandpiper/pkg/shared/rbac"

	al "autocare.org/sandpiper/pkg/api/auth/logging"
	at "autocare.org/sandpiper/pkg/api/auth/transport"
)

// Register ties the auth service to its logger and transport mechanisms
func Register(db *pg.DB, sec auth.Securer, log sandpiper.Logger, srv *echo.Echo, tgen auth.TokenGenerator, mwFunc echo.MiddlewareFunc) {
	svc := auth.Initialize(db, tgen, sec, rbac.New())
	ls := al.ServiceLogger(svc, log)
	at.NewHTTP(ls, srv, mwFunc)
}
