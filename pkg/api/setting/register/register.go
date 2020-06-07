// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package setting

import (
	"github.com/labstack/echo/v4"

	"github.com/sandpiper-framework/sandpiper/pkg/api/setting"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/model"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/rbac"

	gl "github.com/sandpiper-framework/sandpiper/pkg/api/setting/logging"
	gt "github.com/sandpiper-framework/sandpiper/pkg/api/setting/transport"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/database"
)

// Register ties the setting service to its logger and transport mechanisms
func Register(db *database.DB, sec setting.Securer, log sandpiper.Logger, v1 *echo.Group) {
	svc := setting.Initialize(db, rbac.New(db.Settings.ServerRole), sec)
	ls := gl.ServiceLogger(svc, log)
	gt.NewHTTP(ls, v1)
}
