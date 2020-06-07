// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

// Package server sets up the echo server with middleware, binder and validater
// from configuration file allows starting (with graceful shutdown)
package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/sandpiper-framework/sandpiper/pkg/shared/middleware/secure"
)

// New instantiates new Echo server.
func New() *echo.Echo {
	e := echo.New()
	e.Use(
		middleware.Logger(),
		middleware.Recover(),
		secure.CORS(),
		secure.Headers(),
	)
	e.GET("/check", healthCheck)
	e.Validator = &CustomValidator{V: validator.New()}
	custErr := &customErrHandler{e: e}
	e.HTTPErrorHandler = custErr.handler
	e.Binder = &CustomBinder{b: &echo.DefaultBinder{}}
	return e
}

// Settings represents server specific configuration
type Settings struct {
	Port                string
	ReadTimeoutSeconds  int
	WriteTimeoutSeconds int
	ServerRole          string
	ServerID            uuid.UUID
	Debug               bool
}

// Start starts echo server.
func Start(srv *echo.Echo, cfg *Settings) {
	httpServer := &http.Server{
		Addr:         formatPort(cfg.Port),
		ReadTimeout:  time.Duration(cfg.ReadTimeoutSeconds) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeoutSeconds) * time.Second,
	}
	srv.Debug = cfg.Debug
	srv.HideBanner = true

	if srv.Debug {
		srv.GET("/routes", listRoutes)
	}

	fmt.Printf("Server Role: \"%s\"\nServer ID: %s\n\n", cfg.ServerRole, cfg.ServerID)

	// Start server
	go func() {
		if err := srv.StartServer(httpServer); err != nil {
			srv.Logger.Info("Shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		srv.Logger.Fatal(err)
	}
}

func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, "Sandpiper API OK")
}

// listRoutes uses an echo built-in function to return all registered routes.
func listRoutes(c echo.Context) error {
	r, _ := json.Marshal(c.Echo().Routes())
	return c.String(http.StatusOK, string(r))
}

func formatPort(port string) string {
	if port[0:1] != ":" {
		return ":" + port
	}
	return port
}
