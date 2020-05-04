// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package api creates each service used by the server (grouped by api version)
// with middleware, logging and routing, and then starts the web server.
package api

import (
	"errors"

	"sandpiper/pkg/shared/config"
	"sandpiper/pkg/shared/database"
	"sandpiper/pkg/shared/middleware/jwt"
	"sandpiper/pkg/shared/secure"
	"sandpiper/pkg/shared/server"
	"sandpiper/pkg/shared/zlog"

	// One import for each service to register (with identifying alias).
	// Must use a register subdirectory to avoid "import cycle" errors.
	ac "sandpiper/pkg/api/activity/register"
	au "sandpiper/pkg/api/auth/register"
	co "sandpiper/pkg/api/company/register"
	gr "sandpiper/pkg/api/grain/register"
	pa "sandpiper/pkg/api/password/register"
	sl "sandpiper/pkg/api/slice/register"
	su "sandpiper/pkg/api/subscription/register"
	sy "sandpiper/pkg/api/sync/register"
	ta "sandpiper/pkg/api/tag/register"
	us "sandpiper/pkg/api/user/register"
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
	sec := secure.New(cfg.App.MinPasswordStr, cfg.Server.APIKeySecretCode())
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
