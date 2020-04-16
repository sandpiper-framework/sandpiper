// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package sync

import (
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/api/sync"
	"autocare.org/sandpiper/pkg/shared/model"
	"autocare.org/sandpiper/pkg/shared/rbac"

	sl "autocare.org/sandpiper/pkg/api/sync/logging"
	st "autocare.org/sandpiper/pkg/api/sync/transport"
	"autocare.org/sandpiper/pkg/shared/database"
)

// Register ties the sync service to its logger and transport mechanisms
func Register(db *database.DB, sec sync.Securer, log sandpiper.Logger, v1 *echo.Group) {
	svc := sync.Initialize(db, rbac.New(), sec)
	ls := sl.ServiceLogger(svc, log)
	st.NewHTTP(ls, v1)
}
