// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package slice

import (
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/primary/slice"
	"autocare.org/sandpiper/pkg/shared/model"
	"autocare.org/sandpiper/pkg/shared/rbac"

	sl "autocare.org/sandpiper/pkg/primary/slice/logging"
	st "autocare.org/sandpiper/pkg/primary/slice/transport"
)

// Register ties the slice service to its logger and transport mechanisms
func Register(db *pg.DB, sec slice.Securer, log sandpiper.Logger, v1 *echo.Group) {
	svc := slice.Initialize(db, rbac.New(), sec)
	ls := sl.ServiceLogger(svc, log)
	st.NewHTTP(ls, v1)
}
