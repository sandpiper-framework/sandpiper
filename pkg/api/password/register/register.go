// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package password

import (
	"github.com/labstack/echo/v4"

	"github.com/sandpiper-framework/sandpiper/pkg/api/password"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/model"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/rbac"

	pl "github.com/sandpiper-framework/sandpiper/pkg/api/password/logging"
	pt "github.com/sandpiper-framework/sandpiper/pkg/api/password/transport"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/database"
)

// Register ties the company service to its logger and transport mechanisms
func Register(db *database.DB, sec password.Securer, log sandpiper.Logger, v1 *echo.Group) {
	svc := password.Initialize(db, rbac.New(db.Settings.ServerRole), sec)
	ls := pl.ServiceLogger(svc, log)
	pt.NewHTTP(ls, v1)
}
