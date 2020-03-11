// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package user

import (
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/primary/user"
	"autocare.org/sandpiper/pkg/shared/model"
	"autocare.org/sandpiper/pkg/shared/rbac"

	ul "autocare.org/sandpiper/pkg/primary/user/logging"
	ut "autocare.org/sandpiper/pkg/primary/user/transport"
)

// Register ties the user service to its logger and transport mechanisms
func Register(db *pg.DB, sec user.Securer, log sandpiper.Logger, v1 *echo.Group) {
	svc := user.Initialize(db, rbac.New(), sec)
	ls := ul.ServiceLogger(svc, log)
	ut.NewHTTP(ls, v1)
}
