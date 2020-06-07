// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

// Package api creates each service used by the server (grouped by api version)
// with middleware, logging and routing, and then starts the web server.
package api

import (
	"net/http"

	"github.com/GeertJohan/go.rice"
	"github.com/labstack/echo/v4"

	"github.com/sandpiper-framework/sandpiper/pkg/shared/config"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/database"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/middleware/jwt"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/secure"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/server"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/zlog"

	// One import for each service to register (with identifying alias).
	// Must use a register subdirectory to avoid "import cycle" errors.
	ac "github.com/sandpiper-framework/sandpiper/pkg/api/activity/register"
	au "github.com/sandpiper-framework/sandpiper/pkg/api/auth/register"
	co "github.com/sandpiper-framework/sandpiper/pkg/api/company/register"
	gr "github.com/sandpiper-framework/sandpiper/pkg/api/grain/register"
	pa "github.com/sandpiper-framework/sandpiper/pkg/api/password/register"
	se "github.com/sandpiper-framework/sandpiper/pkg/api/setting/register"
	sl "github.com/sandpiper-framework/sandpiper/pkg/api/slice/register"
	su "github.com/sandpiper-framework/sandpiper/pkg/api/subscription/register"
	sy "github.com/sandpiper-framework/sandpiper/pkg/api/sync/register"
	ta "github.com/sandpiper-framework/sandpiper/pkg/api/tag/register"
	us "github.com/sandpiper-framework/sandpiper/pkg/api/user/register"
)

// Start configures and launches the API services
func Start(cfg *config.Configuration) error {

	// setup database connection (with optional query logging using standard "log")
	db, err := database.New(cfg.DB, cfg.DB.Timeout, cfg.DB.LogQueries)
	if err != nil {
		return err
	}

	// setup token, security and logging available for all services
	sec := secure.New(cfg.App.MinPasswordStr, cfg.Server.APIKeySecretCode())
	tok, err := jwt.New(cfg.JWT.SecretKey(), cfg.JWT.SigningAlgorithm, cfg.JWT.Duration, cfg.JWT.MinSecretLength)
	if err != nil {
		return err
	}
	log := zlog.New(cfg.App.ServiceLogging)

	// setup echo server (singleton)
	srv := server.New()

	// sign-up screen (files served from packaged code in `pkged.go`)
	// publicDir := pkger.Dir("/public")
	// relative to this source file
	publicDir := rice.MustFindBox("../../public").HTTPBox()
	srv.GET("/", echo.WrapHandler(http.FileServer(publicDir)))

	// create version group using token authentication middleware
	v1 := srv.Group("/v1")
	v1.Use(tok.MWFunc())

	// register each service (using proper import alias)
	au.Register(db, sec, log, srv, tok, tok.MWFunc()) // auth service (no version group)
	ac.Register(db, sec, log, v1)                     // activity service
	co.Register(db, sec, log, v1)                     // company service
	gr.Register(db, sec, log, v1)                     // grain service
	pa.Register(db, sec, log, v1)                     // password service
	se.Register(db, sec, log, v1)                     // setting service
	sl.Register(db, sec, log, v1)                     // slice service
	su.Register(db, sec, log, v1)                     // subscription service
	sy.Register(db, sec, log, v1)                     // sync (exchange) service
	ta.Register(db, sec, log, v1)                     // tagging service
	us.Register(db, sec, log, v1)                     // user service

	// listen for requests
	server.Start(srv, &server.Settings{
		Port:                cfg.Server.Port,
		ReadTimeoutSeconds:  cfg.Server.ReadTimeout,
		WriteTimeoutSeconds: cfg.Server.WriteTimeout,
		ServerRole:          db.Settings.ServerRole,
		ServerID:            db.Settings.ServerID,
		Debug:               cfg.Server.Debug,
	})

	return nil
}
