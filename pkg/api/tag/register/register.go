// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package tag

import (
	"github.com/labstack/echo/v4"

	"github.com/sandpiper-framework/sandpiper/pkg/api/tag"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/model"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/rbac"

	tl "github.com/sandpiper-framework/sandpiper/pkg/api/tag/logging"
	tt "github.com/sandpiper-framework/sandpiper/pkg/api/tag/transport"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/database"
)

// Register ties the subscription service to its logger and transport mechanisms
func Register(db *database.DB, sec tag.Securer, log sandpiper.Logger, v1 *echo.Group) {
	svc := tag.Initialize(db, rbac.New(db.Settings.ServerRole), sec)
	ls := tl.ServiceLogger(svc, log)
	tt.NewHTTP(ls, v1)
}
