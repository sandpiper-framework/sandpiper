// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package subscription

import (
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/internal/model"

	"autocare.org/sandpiper/pkg/api/subscription"
	sl "autocare.org/sandpiper/pkg/api/subscription/logging"
	st "autocare.org/sandpiper/pkg/api/subscription/transport"
)

// Register ties the subscription service to its logger and transport mechanisms
func Register(db *pg.DB, rbac subscription.RBAC, sec subscription.Securer, log sandpiper.Logger, v1 *echo.Group) {
	svc := subscription.Initialize(db, rbac, sec)
	ls := sl.ServiceLogger(svc, log)
	st.NewHTTP(ls, v1)
}
