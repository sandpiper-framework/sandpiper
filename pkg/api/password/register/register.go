// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package password

import (
	"github.com/labstack/echo/v4"

	"sandpiper/pkg/api/password"
	"sandpiper/pkg/shared/model"
	"sandpiper/pkg/shared/rbac"

	pl "sandpiper/pkg/api/password/logging"
	pt "sandpiper/pkg/api/password/transport"
	"sandpiper/pkg/shared/database"
)

// Register ties the company service to its logger and transport mechanisms
func Register(db *database.DB, sec password.Securer, log sandpiper.Logger, v1 *echo.Group) {
	svc := password.Initialize(db, rbac.New(db.ServerRole), sec)
	ls := pl.ServiceLogger(svc, log)
	pt.NewHTTP(ls, v1)
}
