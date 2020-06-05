// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package tag

import (
	"github.com/labstack/echo/v4"

	"sandpiper/pkg/api/tag"
	"sandpiper/pkg/shared/model"
	"sandpiper/pkg/shared/rbac"

	tl "sandpiper/pkg/api/tag/logging"
	tt "sandpiper/pkg/api/tag/transport"
	"sandpiper/pkg/shared/database"
)

// Register ties the subscription service to its logger and transport mechanisms
func Register(db *database.DB, sec tag.Securer, log sandpiper.Logger, v1 *echo.Group) {
	svc := tag.Initialize(db, rbac.New(db.Settings.ServerRole), sec)
	ls := tl.ServiceLogger(svc, log)
	tt.NewHTTP(ls, v1)
}
