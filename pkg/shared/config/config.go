// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package config is used to read config files and load matching data structures.
package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

const header = `# Copyright The Sandpiper Authors. All rights reserved.
# Use of this source code is governed by an MIT-style
# license that can be found in the LICENSE.md file.

# sandpiper configuration file (rename to "config.yaml" for default use by api server)

`

// Load returns Configuration struct
func Load(path string) (*Configuration, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file \"%s\"", err)
	}
	var cfg = new(Configuration)

	if err := yaml.Unmarshal(b, cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %v", err)
	}
	return cfg, nil
}

// Save creates a config file from a struct to a file
func Save(c *Configuration, filename string) error {
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	buf := bytes.NewBufferString(header)
	if _, err := buf.Write(b); err != nil {
		return err
	}
	return ioutil.WriteFile(filename, buf.Bytes(), 0644)
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
	Dialect    string `yaml:"dialect,omitempty"`
	Network    string `yaml:"network,omitempty"`
	Host       string `yaml:"host,omitempty"`
	Port       string `yaml:"port,omitempty"`
	Database   string `yaml:"database,omitempty"`
	User       string `yaml:"user,omitempty"`
	Password   string `yaml:"password,omitempty"`
	SSLMode    string `yaml:"sslmode,omitempty"`
	Timeout    int    `yaml:"timeout_seconds,omitempty"`
	LogQueries bool   `yaml:"log_queries,omitempty"`
}

// URL creates a connection URL from a `database` section overriding User/Password with env vars if found
func (d *Database) URL() string {
	return d.psn(false)
}

// SafeURL creates a connection URL from a `database` config and env vars without a password
func (d *Database) SafeURL() string {
	return d.psn(true)
}

func (d *Database) psn(safe bool) string {
	password := "*******"
	if !safe {
		password = env("DB_PASSWORD", d.Password)
	}
	if d.Network == "unix" {
		// for unix sockets, use form: postgres://user@:port/database?host=/var/run/postgresql
		return fmt.Sprintf("%s://%s@:%s/%s?host=%s",
			d.Dialect,
			env("DB_USER", d.User),
			env("DB_PORT", d.Port),
			env("DB_DATABASE", d.Database),
			env("DB_HOST", d.Host),
		)
	}
	// postgres://username:password@host:port/database?sslmode=disable
	return fmt.Sprintf("%s://%s:%s@%s/%s?sslmode=%s",
		d.Dialect,
		env("DB_USER", d.User),
		password,
		env("DB_HOST", d.Host),
		env("DB_DATABASE", d.Database),
		env("DB_SSLMODE", d.SSLMode),
	)
}

// DSN provides a connection string using key/value pairs format from on a database config
// https://godoc.org/github.com/lib/pq#hdr-Connection_String_Parameters
func (d *Database) DSN() string {
	b := make([]string,10)

	var add = func(k, v, e string) {
		val := env(e, v)
		if val != "" {
			b = append(b, k + "=" + v)
		}
	}

	add("dbname",d.Database,"DB_DATABASE")
	add("user", d.User, "DB_USER")
	add("password", d.Password, "DB_PASSWORD")
	add("host", d.Host, "DB_HOST")
	add("port", d.Port, "DB_PORT")
	add("sslmode", d.SSLMode, "DB_SSLMODE")

	return strings.Join(b, " ")
}

// Server holds data necessary for server configuration
type Server struct {
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
