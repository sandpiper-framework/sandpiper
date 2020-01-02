// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package api creates each service used by the server (grouped by api version)
// with middleware, logging and routing, and then starts the web server.
package api

import (
	"autocare.org/sandpiper/pkg/internal/database"
	"autocare.org/sandpiper/pkg/internal/middleware/jwt"
	"autocare.org/sandpiper/pkg/internal/rbac"
	"autocare.org/sandpiper/pkg/internal/secure"
	"autocare.org/sandpiper/pkg/internal/server"
	"autocare.org/sandpiper/pkg/internal/zlog"
	"autocare.org/sandpiper/pkg/shared/config"

	// one import for each service to register (with identifying alias)
	ar "autocare.org/sandpiper/pkg/api/auth/register"
	cr "autocare.org/sandpiper/pkg/api/company/register"
	gr "autocare.org/sandpiper/pkg/api/grain/register"
	pr "autocare.org/sandpiper/pkg/api/password/register"
	sr "autocare.org/sandpiper/pkg/api/slice/register"
	ur "autocare.org/sandpiper/pkg/api/user/register"
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
	rba := rbac.New()
	tok := jwt.New(cfg.JWT.Secret, cfg.JWT.SigningAlgorithm, cfg.JWT.Duration)
	log := zlog.New(cfg.App.ServiceLogging)

	// setup echo server (singleton)
	srv := server.New()

	// create version group using token middleware
	v1 := srv.Group("/v1")
	v1.Use(tok.MWFunc())

	// register each service (using proper import alias)
	ar.Register(db, rba, sec, log, srv, tok, tok.MWFunc()) // auth service (no version group)
	pr.Register(db, rba, sec, log, v1)                     // password service
	ur.Register(db, rba, sec, log, v1)                     // user service
	cr.Register(db, rba, sec, log, v1)                     // company service
	sr.Register(db, rba, sec, log, v1)                     // slice service
	gr.Register(db, rba, sec, log, v1)                     // grain service
	// rr.Register(db, rba, sec, log, v1)  // subscription service
	// xr.Register(db, rba, sec, log, v1)  // sync (exchange) service

	// listen for requests
	server.Start(srv, &server.Config{
		Port:                cfg.Server.Port,
		ReadTimeoutSeconds:  cfg.Server.ReadTimeout,
		WriteTimeoutSeconds: cfg.Server.WriteTimeout,
		Debug:               cfg.Server.Debug,
	})

	return nil
}
