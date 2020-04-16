// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package user

import (
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/api/user"
	"autocare.org/sandpiper/pkg/shared/database"
	"autocare.org/sandpiper/pkg/shared/model"
	"autocare.org/sandpiper/pkg/shared/rbac"

	ul "autocare.org/sandpiper/pkg/api/user/logging"
	ut "autocare.org/sandpiper/pkg/api/user/transport"
)

// Register ties the user service to its logger and transport mechanisms
func Register(db *database.DB, sec user.Securer, log sandpiper.Logger, v1 *echo.Group) {
	svc := user.Initialize(db, rbac.New(), sec)
	ls := ul.ServiceLogger(svc, log)
	ut.NewHTTP(ls, v1)
}
