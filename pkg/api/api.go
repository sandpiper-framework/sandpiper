// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// The api package creates each service used by the server (grouped by api version)
// with middleware, logging and routing, and then starts the echo web server.
package api

import (
	"autocare.org/sandpiper/pkg/api/auth"
	al "autocare.org/sandpiper/pkg/api/auth/logging"
	at "autocare.org/sandpiper/pkg/api/auth/transport"
	"crypto/sha1"

	"autocare.org/sandpiper/pkg/api/password"
	pl "autocare.org/sandpiper/pkg/api/password/logging"
	pt "autocare.org/sandpiper/pkg/api/password/transport"

	"autocare.org/sandpiper/pkg/api/user"
	ul "autocare.org/sandpiper/pkg/api/user/logging"
	ut "autocare.org/sandpiper/pkg/api/user/transport"

	"autocare.org/sandpiper/pkg/config"
	"autocare.org/sandpiper/pkg/middleware/jwt"
	"autocare.org/sandpiper/pkg/postgres"
	"autocare.org/sandpiper/pkg/rbac"
	"autocare.org/sandpiper/pkg/secure"
	"autocare.org/sandpiper/pkg/server"
	"autocare.org/sandpiper/pkg/zlog"
)

// Start configures and launches the API services
func Start(cfg *config.Configuration) error {

	// setup database connection with optional query debug logging (using standard "log")
	db, err := postgres.New(cfg.DB.PSN, cfg.DB.Timeout, cfg.DB.LogQueries)
	if err != nil {
		return err
	}

	// setup middleware and logging services
	sec := secure.New(cfg.App.MinPasswordStr, sha1.New())
	rba := rbac.New()
	tok := jwt.New(cfg.JWT.Secret, cfg.JWT.SigningAlgorithm, cfg.JWT.Duration)
	log := zlog.New()

	// setup server
	srv := server.New()
	srv.Static("/swaggerui", cfg.App.SwaggerUIPath)

	// auth service is special (doesn't include api version)
	at.NewHTTP(al.New(auth.Initialize(db, tok, sec, rba), log), srv, tok.MWFunc())

	v1 := srv.Group("/v1")
	v1.Use(tok.MWFunc())

	// user service
	ut.NewHTTP(ul.New(user.Initialize(db, rba, sec), log), v1)

	// password service
	pt.NewHTTP(pl.New(password.Initialize(db, rba, sec), log), v1)

	// kick it off
	server.Start(srv, &server.Config{
		Port:                cfg.Server.Port,
		ReadTimeoutSeconds:  cfg.Server.ReadTimeout,
		WriteTimeoutSeconds: cfg.Server.WriteTimeout,
		Debug:               cfg.Server.Debug,
		HideBanner:          true, // hides echo banner (we have our own)... not in config file
	})

	return nil
}
