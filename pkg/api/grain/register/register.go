// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package grain

import (
	"github.com/labstack/echo/v4"

	"sandpiper/pkg/api/grain"
	"sandpiper/pkg/shared/model"
	"sandpiper/pkg/shared/rbac"

	gl "sandpiper/pkg/api/grain/logging"
	gt "sandpiper/pkg/api/grain/transport"
	"sandpiper/pkg/shared/database"
)

// Register ties the grain service to its logger and transport mechanisms
func Register(db *database.DB, sec grain.Securer, log sandpiper.Logger, v1 *echo.Group) {
	svc := grain.Initialize(db, rbac.New(db.ServerRole), sec)
	ls := gl.ServiceLogger(svc, log)
	gt.NewHTTP(ls, v1)
}
