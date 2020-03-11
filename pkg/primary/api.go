// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package primary creates each service used by the primary server (grouped by api version)
// with middleware, logging and routing, and then starts the web server.
package primary

import (
	"autocare.org/sandpiper/pkg/shared/config"
	"autocare.org/sandpiper/pkg/shared/database"
	"autocare.org/sandpiper/pkg/shared/middleware/jwt"
	"autocare.org/sandpiper/pkg/shared/secure"
	"autocare.org/sandpiper/pkg/shared/server"
	"autocare.org/sandpiper/pkg/shared/zlog"

	// One import for each service to register (with identifying alias).
	// Must use a register subdirectory to avoid "import cycle" errors.
	au "autocare.org/sandpiper/pkg/primary/auth/register"
	co "autocare.org/sandpiper/pkg/primary/company/register"
	gr "autocare.org/sandpiper/pkg/primary/grain/register"
	pa "autocare.org/sandpiper/pkg/primary/password/register"
	sl "autocare.org/sandpiper/pkg/primary/slice/register"
	su "autocare.org/sandpiper/pkg/primary/subscription/register"
	sy "autocare.org/sandpiper/pkg/primary/sync/register"
	ta "autocare.org/sandpiper/pkg/primary/tag/register"
	us "autocare.org/sandpiper/pkg/primary/user/register"
)

// Start configures and launches the API services
func Start(cfg *config.Configuration) error {

	// setup database connection (with optional query logging using standard "log")
	db, err := database.New(cfg.DB.URL(), cfg.DB.Timeout, cfg.DB.LogQueries)
	if err != nil {
		return err
	}

	// setup token, security and logging available for all services
	sec := secure.New(cfg.App.MinPasswordStr)
	tok := jwt.New(cfg.JWT.Secret, cfg.JWT.SigningAlgorithm, cfg.JWT.Duration)
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
	sy.Register(db, sec, log, v1)                     // sync (exchange) service

	// listen for requests
	server.Start(srv, &server.Config{
		Port:                cfg.Server.Port,
		ReadTimeoutSeconds:  cfg.Server.ReadTimeout,
		WriteTimeoutSeconds: cfg.Server.WriteTimeout,
		Debug:               cfg.Server.Debug,
	})

	return nil
}
