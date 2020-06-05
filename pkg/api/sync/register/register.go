// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package sync

import (
	"github.com/labstack/echo/v4"
	"sandpiper/pkg/api/sync"
	"sandpiper/pkg/shared/database"
	"sandpiper/pkg/shared/model"
	"sandpiper/pkg/shared/rbac"

	sl "sandpiper/pkg/api/sync/logging"
	st "sandpiper/pkg/api/sync/transport"
)

// Register ties the sync service to its logger and transport mechanisms
func Register(db *database.DB, sec sync.Securer, log sandpiper.Logger, v1 *echo.Group) {
	rba := rbac.New(db.Settings.ServerRole)
	rba.ServerID = db.Settings.ServerID
	svc := sync.Initialize(db, rba, sec)
	ls := sl.ServiceLogger(svc, log)
	st.NewHTTP(ls, v1)
}
