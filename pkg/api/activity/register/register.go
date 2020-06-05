// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package sync

import (
	"github.com/labstack/echo/v4"

	"sandpiper/pkg/api/activity"
	"sandpiper/pkg/shared/model"
	"sandpiper/pkg/shared/rbac"

	sl "sandpiper/pkg/api/activity/logging"
	st "sandpiper/pkg/api/activity/transport"
	"sandpiper/pkg/shared/database"
)

// Register ties the sync service to its logger and transport mechanisms
func Register(db *database.DB, sec activity.Securer, log sandpiper.Logger, v1 *echo.Group) {
	svc := activity.Initialize(db, rbac.New(db.Settings.ServerRole), sec)
	ls := sl.ServiceLogger(svc, log)
	st.NewHTTP(ls, v1)
}
