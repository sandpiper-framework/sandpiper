// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package api

import (
	"crypto/sha1"

	"autocare.org/sandpiper/pkg/api/auth"
	al "autocare.org/sandpiper/pkg/api/auth/logging"
	at "autocare.org/sandpiper/pkg/api/auth/transport"

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

// Start starts the API service
func Start(cfg *config.Configuration) error {
	db, err := postgres.New(cfg.DB.PSN, cfg.DB.Timeout, cfg.DB.LogQueries)
	if err != nil {
		return err
	}

	sec := secure.New(cfg.App.MinPasswordStr, sha1.New())
	rba := rbac.New()
	tok := jwt.New(cfg.JWT.Secret, cfg.JWT.SigningAlgorithm, cfg.JWT.Duration)
	log := zlog.New()

	srv := server.New()
	srv.Static("/swaggerui", cfg.App.SwaggerUIPath)

	at.NewHTTP(al.New(auth.Initialize(db, tok, sec, rba), log), srv, tok.MWFunc())

	v1 := srv.Group("/v1")
	v1.Use(tok.MWFunc())

	ut.NewHTTP(ul.New(user.Initialize(db, rba, sec), log), v1)
	pt.NewHTTP(pl.New(password.Initialize(db, rba, sec), log), v1)

	server.Start(srv, &server.Config{
		Port:                cfg.Server.Port,
		ReadTimeoutSeconds:  cfg.Server.ReadTimeout,
		WriteTimeoutSeconds: cfg.Server.WriteTimeout,
		Debug:               cfg.Server.Debug,
	})

	return nil
}
