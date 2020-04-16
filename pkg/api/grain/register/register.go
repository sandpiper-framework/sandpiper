// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package grain

import (
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/api/grain"
	"autocare.org/sandpiper/pkg/shared/model"
	"autocare.org/sandpiper/pkg/shared/rbac"

	gl "autocare.org/sandpiper/pkg/api/grain/logging"
	gt "autocare.org/sandpiper/pkg/api/grain/transport"
	"autocare.org/sandpiper/pkg/shared/database"
)

// Register ties the grain service to its logger and transport mechanisms
func Register(db *database.DB, sec grain.Securer, log sandpiper.Logger, v1 *echo.Group) {
	svc := grain.Initialize(db, rbac.New(), sec)
	ls := gl.ServiceLogger(svc, log)
	gt.NewHTTP(ls, v1)
}
