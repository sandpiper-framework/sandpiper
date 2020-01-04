package auth

import (
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/api/auth"
	al "autocare.org/sandpiper/pkg/api/auth/logging"
	at "autocare.org/sandpiper/pkg/api/auth/transport"
	"autocare.org/sandpiper/pkg/internal/model"
)

// Register ties the company service to its logger and transport mechanisms
func Register(db *pg.DB, rbac auth.RBAC, sec auth.Securer, log sandpiper.Logger, srv *echo.Echo, tgen auth.TokenGenerator, mwFunc echo.MiddlewareFunc) {
	svc := auth.Initialize(db, tgen, sec, rbac)
	ls := al.ServiceLogger(svc, log)
	at.NewHTTP(ls, srv, mwFunc)
}
