// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package config is used to read config files and load matching data structures.
package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// Load returns Configuration struct
func Load(path string) (*Configuration, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file \"%s\"", err)
	}
	var cfg = new(Configuration)

	if err := yaml.Unmarshal(bytes, cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %v", err)
	}
	return cfg, nil
}

// Configuration defines available config sections with pointers to their structs
type Configuration struct {
	Server  *Server      `yaml:"server,omitempty"`
	DB      *Database    `yaml:"database,omitempty"`
	JWT     *JWT         `yaml:"jwt,omitempty"`
	App     *Application `yaml:"application,omitempty"`
	Command *Command     `yaml:"command,omitempty"`
}

// Database structure holds settings for database configuration
type Database struct {
	LogQueries bool   `yaml:"log_queries,omitempty"`
	Timeout    int    `yaml:"timeout_seconds,omitempty"`
	Dialect    string `yaml:"dialect,omitempty"`
	Database   string `yaml:"database,omitempty"`
	User       string `yaml:"user,omitempty"`     // can override from env var
	Password   string `yaml:"password,omitempty"` // can override from env var
	Host       string `yaml:"host,omitempty"`
	Port       string `yaml:"port,omitempty"`
	SSLMode    string `yaml:"sslmode,omitempty"`
}

// URL creates a connection URL from a `database` section overriding User/Password with env vars if found
func (d *Database) URL() string {
	// postgres://username:password@host:port/database?sslmode=disable
	return fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s",
		d.Dialect,
		env("DB_USER", d.User),
		env("DB_PASSWORD", d.Password),
		d.Host,
		d.Port,
		d.Database,
		d.SSLMode)
}

// SafeURL creates a connection URL from a `database` config and env vars without a password
func (d *Database) SafeURL() string {
	// postgres://username@host:port/database?sslmode=disable
	return fmt.Sprintf("%s://%s@%s:%s/%s?sslmode=%s",
		d.Dialect,
		env("DB_USER", d.User),
		d.Host,
		d.Port,
		d.Database,
		d.SSLMode)
}

// Server holds data necessary for server configuration
type Server struct {
	URL          string `yaml:"url,omitempty"` // is this being used?
	Port         string `yaml:"port,omitempty"`
	Debug        bool   `yaml:"debug,omitempty"`
	ReadTimeout  int    `yaml:"read_timeout_seconds,omitempty"`
	WriteTimeout int    `yaml:"write_timeout_seconds,omitempty"`
	MaxSyncProcs int    `yaml:"sync_pool,omitempty"`
	APIKeySecret string `yaml:"api_key_secret,omitempty"`
}

// APIKeySecretCode allows overriding the config value with APIKEY_SECRET environment variable
func (s *Server) APIKeySecretCode() string {
	return env("APIKEY_SECRET", s.APIKeySecret)
}

// JWT holds data necessary for JWT configuration
type JWT struct {
	Secret           string `yaml:"secret,omitempty"`
	Duration         int    `yaml:"duration_minutes,omitempty"`
	RefreshDuration  int    `yaml:"refresh_duration_minutes,omitempty"`
	MaxRefresh       int    `yaml:"max_refresh_minutes,omitempty"`
	SigningAlgorithm string `yaml:"signing_algorithm,omitempty"`
	MinSecretLength  int    `yaml:"min_secret_length,omitempty"`
}

// SecretKey allows overriding the config secret with the JWT_SECRET environment variable
func (j *JWT) SecretKey() string {
	return env("JWT_SECRET", j.Secret)
}

// Application holds application configuration details
type Application struct {
	MinPasswordStr int  `yaml:"min_password_strength,omitempty"`
	ServiceLogging bool `yaml:"service_logging,omitempty"`
}

// Command holds configuration options for the `sandpiper` command
type Command struct {
	URL          string `yaml:"url,omitempty"`
	Port         string `yaml:"port,omitempty"`
	MaxSyncProcs int    `yaml:"max_sync_procs,omitempty"`
	Debug        bool   `yaml:"debug,omitempty"`
}

func env(key, defValue string) string {
	envValue := os.Getenv(key)
	if envValue != "" {
		return envValue
	}
	return defValue
}
