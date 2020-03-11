// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package subscription

import (
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/primary/subscription"
	"autocare.org/sandpiper/pkg/shared/model"
	"autocare.org/sandpiper/pkg/shared/rbac"

	sl "autocare.org/sandpiper/pkg/primary/subscription/logging"
	st "autocare.org/sandpiper/pkg/primary/subscription/transport"
)

// Register ties the subscription service to its logger and transport mechanisms
func Register(db *pg.DB, sec subscription.Securer, log sandpiper.Logger, v1 *echo.Group) {
	svc := subscription.Initialize(db, rbac.New(), sec)
	ls := sl.ServiceLogger(svc, log)
	st.NewHTTP(ls, v1)
}
