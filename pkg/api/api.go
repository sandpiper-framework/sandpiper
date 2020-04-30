// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package api creates each service used by the server (grouped by api version)
// with middleware, logging and routing, and then starts the web server.
package api

import (
	"errors"

	"autocare.org/sandpiper/pkg/shared/config"
	"autocare.org/sandpiper/pkg/shared/database"
	"autocare.org/sandpiper/pkg/shared/middleware/jwt"
	"autocare.org/sandpiper/pkg/shared/secure"
	"autocare.org/sandpiper/pkg/shared/server"
	"autocare.org/sandpiper/pkg/shared/zlog"

	// One import for each service to register (with identifying alias).
	// Must use a register subdirectory to avoid "import cycle" errors.
	ac "autocare.org/sandpiper/pkg/api/activity/register"
	au "autocare.org/sandpiper/pkg/api/auth/register"
	co "autocare.org/sandpiper/pkg/api/company/register"
	gr "autocare.org/sandpiper/pkg/api/grain/register"
	pa "autocare.org/sandpiper/pkg/api/password/register"
	sl "autocare.org/sandpiper/pkg/api/slice/register"
	su "autocare.org/sandpiper/pkg/api/subscription/register"
	sy "autocare.org/sandpiper/pkg/api/sync/register"
	ta "autocare.org/sandpiper/pkg/api/tag/register"
	us "autocare.org/sandpiper/pkg/api/user/register"
)

const (
	// ServerRoleKey accesses the server role setting ("primary" or "secondary")
	ServerRoleKey = "server-role"
)

// Start configures and launches the API services
func Start(cfg *config.Configuration) error {

	// setup database connection (with optional query logging using standard "log")
	db, err := database.New(cfg.DB.URL(), cfg.DB.Timeout, cfg.DB.LogQueries)
	if err != nil {
		return err
	}

	// check required db setting
	db.ServerRole = db.Settings[ServerRoleKey]
	if db.ServerRole == "" {
		return errors.New("missing db setting: \"" + ServerRoleKey + "\"")
	}

	// setup token, security and logging available for all services
	sec := secure.New(cfg.App.MinPasswordStr, cfg.Server.APIKeySecret)
	tok, err := jwt.New(cfg.JWT.SecretKey(), cfg.JWT.SigningAlgorithm, cfg.JWT.Duration, cfg.JWT.MinSecretLength)
	if err != nil {
		return err
	}
	log := zlog.New(cfg.App.ServiceLogging)

	// setup echo server (singleton)
	srv := server.New()

	// create version group using token authentication middleware
	v1 := srv.Group("/v1")
	v1.Use(tok.MWFunc())

	// register each service (using proper import alias)
	au.Register(db, sec, log, srv, tok, tok.MWFunc()) // auth service (no version group)
	pa.Register(db, sec, log, v1)                     // password service
	us.Register(db, sec, log, v1)                     // user service
	co.Register(db, sec, log, v1)                     // company service
	sl.Register(db, sec, log, v1)                     // slice service
	gr.Register(db, sec, log, v1)                     // grain service
	su.Register(db, sec, log, v1)                     // subscription service
	ta.Register(db, sec, log, v1)                     // tagging service
	ac.Register(db, sec, log, v1)                     // activity service
	sy.Register(db, sec, log, v1)                     // sync (exchange) service

	// listen for requests
	server.Start(srv, &server.Settings{
		Port:                cfg.Server.Port,
		ReadTimeoutSeconds:  cfg.Server.ReadTimeout,
		WriteTimeoutSeconds: cfg.Server.WriteTimeout,
		ServerRole:          db.ServerRole,
		Debug:               cfg.Server.Debug,
	})

	return nil
}
