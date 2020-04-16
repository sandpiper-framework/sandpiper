// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package tag

import (
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/api/tag"
	"autocare.org/sandpiper/pkg/shared/model"
	"autocare.org/sandpiper/pkg/shared/rbac"

	tl "autocare.org/sandpiper/pkg/api/tag/logging"
	tt "autocare.org/sandpiper/pkg/api/tag/transport"
	"autocare.org/sandpiper/pkg/shared/database"
)

// Register ties the subscription service to its logger and transport mechanisms
func Register(db *database.DB, sec tag.Securer, log sandpiper.Logger, v1 *echo.Group) {
	svc := tag.Initialize(db, rbac.New(), sec)
	ls := tl.ServiceLogger(svc, log)
	tt.NewHTTP(ls, v1)
}
