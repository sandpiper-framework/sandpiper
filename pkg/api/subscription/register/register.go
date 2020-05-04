// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package subscription

import (
	"github.com/labstack/echo/v4"

	"sandpiper/pkg/api/subscription"
	"sandpiper/pkg/shared/model"
	"sandpiper/pkg/shared/rbac"

	sl "sandpiper/pkg/api/subscription/logging"
	st "sandpiper/pkg/api/subscription/transport"
	"sandpiper/pkg/shared/database"
)

// Register ties the subscription service to its logger and transport mechanisms
func Register(db *database.DB, sec subscription.Securer, log sandpiper.Logger, v1 *echo.Group) {
	svc := subscription.Initialize(db, rbac.New(), sec)
	ls := sl.ServiceLogger(svc, log)
	st.NewHTTP(ls, v1)
}
