// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package subscription

import (
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/api/subscription"
	"autocare.org/sandpiper/pkg/shared/model"
	"autocare.org/sandpiper/pkg/shared/rbac"

	sl "autocare.org/sandpiper/pkg/api/subscription/logging"
	st "autocare.org/sandpiper/pkg/api/subscription/transport"
	"autocare.org/sandpiper/pkg/shared/database"
)

// Register ties the subscription service to its logger and transport mechanisms
func Register(db *database.DB, sec subscription.Securer, log sandpiper.Logger, v1 *echo.Group) {
	svc := subscription.Initialize(db, rbac.New(), sec)
	ls := sl.ServiceLogger(svc, log)
	st.NewHTTP(ls, v1)
}