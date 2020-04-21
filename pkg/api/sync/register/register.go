// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package sync

import (
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/api/sync"
	"autocare.org/sandpiper/pkg/shared/database"
	"autocare.org/sandpiper/pkg/shared/model"
	"autocare.org/sandpiper/pkg/shared/rbac"

	sl "autocare.org/sandpiper/pkg/api/sync/logging"
	st "autocare.org/sandpiper/pkg/api/sync/transport"
)

// Register ties the sync service to its logger and transport mechanisms
func Register(db *database.DB, sec sync.Securer, log sandpiper.Logger, v1 *echo.Group) {
	rba := rbac.New()
	rba.ServerRole = db.ServerRole // save server-role for rba checks
	svc := sync.Initialize(db, rba, sec)
	ls := sl.ServiceLogger(svc, log)
	st.NewHTTP(ls, v1)
}
