// mockapi creates an http server that mimics responses to a subset of the
// sandpiper api for client testing purposes (without a database and before the
// sandpiper server is completed). We anticipate this utility having a short lifespan.

package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
)

// Config represents server specific config.
type Config struct {
	Port                string
	ReadTimeoutSeconds  int
	WriteTimeoutSeconds int
	Debug               bool
}

func main() {

	// todo: pull config from external yaml file (or maybe from args)
	cfg := Config{
		Port:                "localhost:3030",
		ReadTimeoutSeconds:  5,
		WriteTimeoutSeconds: 5,
		Debug:               true,
	}

	// instantiate http server
	srv := echo.New()

	// configure routes
	v1 := srv.Group("/v1")
	v1.GET("/routes", listRoutes)
	v1.GET("/login", login)
	v1.GET("/slices", getMySlices)
	//v1.GET("/slices/:id", getSliceInfo)
	v1.POST("/slices/:id", postObject)

	// kick it off
	Start(srv, &cfg)
}

// postObject adds a new object to the slice.
func postObject(c echo.Context) error {
	// todo: provide error for slice not found condition
	// todo: provide error for duplicate object-id
	return c.String(http.StatusOK, "OK")
}

// getMySlices only returns slice information assigned to the login user.
func getMySlices(c echo.Context) error {
	type Slice struct {
		SliceID   string `json:"slice-id"`
		SliceName string `json:"slice-name"`
		SliceHash string `json:"slice-hash"`
	}

	type resp struct {
		Slices []Slice `json:"slices"`
		Count  int     `json:"count"`
	}

	slice := []Slice{
		{"08efdf90-a815-4cf7-b71c-008e5fd31cce", "AAP-Brakes", "cf23df2207d99a74fbe169e3eba035e633b65d94"},
		{"cb4b768b-6d6b-4965-a29a-9052a80dbbbb", "AAP-Wipers", "1a804c61e1a70ab37b912792ee846de7378c4a36"},
	}

	var r = resp{
		Slices: slice,
		Count:  len(slice),
	}
	return c.JSON(http.StatusOK, r)
}

// login returns the token which should be used in all other calls. It does not currently
// validate the request data (username and password).
func login(c echo.Context) error {
	const (
		tok = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9" +
			".eyJjIjoxLCJlIjoiam9obmRvZUBtYWlsLmNvbSIsImV4cCI6MTU3NjI3Mzg4NywiaWQiOjEsImwiOjEsInIiOjEwMCwidSI6ImFkbWluIn0" +
			".0u9B4GjwmI2VnEhEWVdp5khxgvNoHPR8yGj2f3n4PKY"
	)
	type resp struct {
		Token      string `json:"token"`
		Expires    string `json:"expires"`
		RefreshTok string `json:"refresh_token"`
	}
	var r = resp{
		Token:      tok,
		Expires:    time.Now().Add(time.Hour).Format(time.RFC3339),
		RefreshTok: "04782f813406b7686fc83f7aa43e694d2b3b9004",
	}
	return c.JSON(http.StatusOK, r)
}

// listRoutes uses an echo built-in function to return all registered routes.
func listRoutes(c echo.Context) error {
	r, _ := json.Marshal(c.Echo().Routes())
	return c.String(http.StatusOK, string(r))
}

// Start starts the echo server in a separate goroutine.
func Start(srv *echo.Echo, cfg *Config) {
	httpServer := &http.Server{
		Addr:         cfg.Port,
		ReadTimeout:  time.Duration(cfg.ReadTimeoutSeconds) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeoutSeconds) * time.Second,
	}
	srv.Debug = cfg.Debug
	srv.HideBanner = true

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
