// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package password

import (
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/primary/password"
	"autocare.org/sandpiper/pkg/shared/model"
	"autocare.org/sandpiper/pkg/shared/rbac"

	pl "autocare.org/sandpiper/pkg/primary/password/logging"
	pt "autocare.org/sandpiper/pkg/primary/password/transport"
)

// Register ties the company service to its logger and transport mechanisms
func Register(db *pg.DB, sec password.Securer, log sandpiper.Logger, v1 *echo.Group) {
	svc := password.Initialize(db, rbac.New(), sec)
	ls := pl.ServiceLogger(svc, log)
	pt.NewHTTP(ls, v1)
}
