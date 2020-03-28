// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"autocare.org/sandpiper/pkg/shared/config"
)

func TestLoadConfig(t *testing.T) {
	cases := []struct {
		name     string
		path     string
		wantData *config.Configuration
		wantErr  bool
	}{
		{
			name:    "Fail on non-existing file",
			path:    "notExists",
			wantErr: true,
		},
		{
			name:    "Fail on wrong file format",
			path:    "testdata/config.invalid.yaml",
			wantErr: true,
		},
		{
			name: "Match from config",
			path: "testdata/config.testdata.yaml",
			wantData: &config.Configuration{
				DB: &config.Database{
					LogQueries: true,
					Timeout:    5,
					Dialect:    "postgres",
					Database:   "sandpiper",
					User:       "admin",
					Password:   "secret",
					Host:       "localhost",
					Port:       "1234",
					SSLMode:    "disable",
				},
				Server: &config.Server{
					Port:         ":8080",
					Debug:        true,
					ReadTimeout:  10,
					WriteTimeout: 5,
				},
				JWT: &config.JWT{
					Secret:           "jwtrealm",
					Duration:         15,
					RefreshDuration:  20,
					MaxRefresh:       1440,
					SigningAlgorithm: "HS256",
				},
				App: &config.Application{
					MinPasswordStr: 3,
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := config.Load(tt.path)
			assert.Equal(t, tt.wantData, cfg)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
