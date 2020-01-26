// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package tag

import (
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/api/tag"
	"autocare.org/sandpiper/pkg/shared/model"
	"autocare.org/sandpiper/pkg/shared/rbac"

	tl "autocare.org/sandpiper/pkg/api/tag/logging"
	tt "autocare.org/sandpiper/pkg/api/tag/transport"
)

// Register ties the subscription service to its logger and transport mechanisms
func Register(db *pg.DB, sec tag.Securer, log sandpiper.Logger, v1 *echo.Group) {
	svc := tag.Initialize(db, rbac.New(), sec)
	ls := tl.ServiceLogger(svc, log)
	tt.NewHTTP(ls, v1)
}
