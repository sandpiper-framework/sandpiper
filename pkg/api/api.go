// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package api creates each service used by the server (grouped by api version)
// with middleware, logging and routing, and then starts the web server.
package api

import (
	"autocare.org/sandpiper/pkg/api/auth"
	al "autocare.org/sandpiper/pkg/api/auth/logging"
	at "autocare.org/sandpiper/pkg/api/auth/transport"

	"autocare.org/sandpiper/pkg/api/password"
	pl "autocare.org/sandpiper/pkg/api/password/logging"
	pt "autocare.org/sandpiper/pkg/api/password/transport"

	"autocare.org/sandpiper/pkg/api/user"
	ul "autocare.org/sandpiper/pkg/api/user/logging"
	ut "autocare.org/sandpiper/pkg/api/user/transport"

	"autocare.org/sandpiper/pkg/api/company"
	cl "autocare.org/sandpiper/pkg/api/company/logging"
	ct "autocare.org/sandpiper/pkg/api/company/transport"

	"autocare.org/sandpiper/pkg/api/slice"
	sl "autocare.org/sandpiper/pkg/api/slice/logging"
	st "autocare.org/sandpiper/pkg/api/slice/transport"

	"autocare.org/sandpiper/internal/config"
	"autocare.org/sandpiper/internal/database"
	"autocare.org/sandpiper/internal/middleware/jwt"
	"autocare.org/sandpiper/internal/rbac"
	"autocare.org/sandpiper/internal/secure"
	"autocare.org/sandpiper/internal/server"
	"autocare.org/sandpiper/internal/zlog"
)

// Start configures and launches the API services
func Start(cfg *config.Configuration) error {

	// setup database connection with optional query logging (using standard "log")
	db, err := database.New(cfg.DB.URL(), cfg.DB.Timeout, cfg.DB.LogQueries)
	if err != nil {
		return err
	}

	// setup middleware and logging services
	sec := secure.New(cfg.App.MinPasswordStr)
	rba := rbac.New()
	tok := jwt.New(cfg.JWT.Secret, cfg.JWT.SigningAlgorithm, cfg.JWT.Duration)
	log := zlog.New()

	// setup echo server (singleton)
	srv := server.New()

	// auth service is special (doesn't include api version)
	at.NewHTTP(al.New(auth.Initialize(db, tok, sec, rba), log), srv, tok.MWFunc())

	v1 := srv.Group("/v1")
	v1.Use(tok.MWFunc())

	// register each service (using proper import aliases)
	ut.NewHTTP(ul.New(user.Initialize(db, rba, sec), log), v1)     // user service
	pt.NewHTTP(pl.New(password.Initialize(db, rba, sec), log), v1) // password service
	ct.NewHTTP(cl.New(company.Initialize(db, rba, sec), log), v1)  // company service
	st.NewHTTP(sl.New(slice.Initialize(db, rba, sec), log), v1)    // slice service

	// start the server listening
	server.Start(srv, &server.Config{
		Port:                cfg.Server.Port,
		ReadTimeoutSeconds:  cfg.Server.ReadTimeout,
		WriteTimeoutSeconds: cfg.Server.WriteTimeout,
		Debug:               cfg.Server.Debug,
	})

	return nil
}
