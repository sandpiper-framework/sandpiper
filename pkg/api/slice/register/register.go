// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package slice

import (
	"github.com/labstack/echo/v4"

	"sandpiper/pkg/api/slice"
	"sandpiper/pkg/shared/model"
	"sandpiper/pkg/shared/rbac"

	sl "sandpiper/pkg/api/slice/logging"
	st "sandpiper/pkg/api/slice/transport"
	"sandpiper/pkg/shared/database"
)

// Register ties the slice service to its logger and transport mechanisms
func Register(db *database.DB, sec slice.Securer, log sandpiper.Logger, v1 *echo.Group) {
	svc := slice.Initialize(db, rbac.New(db.Settings.ServerRole), sec)
	ls := sl.ServiceLogger(svc, log)
	st.NewHTTP(ls, v1)
}
