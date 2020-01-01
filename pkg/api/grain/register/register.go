// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package grain

import (
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/internal/model"

	"autocare.org/sandpiper/pkg/api/grain"
	gl "autocare.org/sandpiper/pkg/api/grain/logging"
	gt "autocare.org/sandpiper/pkg/api/grain/transport"
)

// Register ties the company service to its logger and transport mechanisms
func Register(db *pg.DB, rbac grain.RBAC, sec grain.Securer, log sandpiper.Logger, v1 *echo.Group) {
	svc := grain.Initialize(db, rbac, sec)
  ls := gl.ServiceLogger(svc, log)
	gt.NewHTTP(ls, v1)
}