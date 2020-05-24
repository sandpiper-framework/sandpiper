// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package database

import (
	"database/sql"
	"fmt"

	"github.com/GuiaBolso/darwin"
	_ "github.com/lib/pq"
)

// defineSchema returns a list of our database migrations
// Each migration script is defined in a separate descriptive string variable.
func defineSchema() []darwin.Migration {
	var (
		enums = `
		CREATE TYPE server_role_enum AS ENUM (
			'primary',
			'secondary'
		);
		CREATE TYPE slice_type_enum AS ENUM (
			'aces-file',
			'aces-items',
			'asset-files',
			'pies-file',
			'pies-items',
			'pies-marketcopy',
			'pies-pricesheet',
			'partspro-file'
		);
		CREATE TYPE encoding_enum AS ENUM (
			'raw',
			'b64',
			'z64',
			'a85',
			'z85'
		);`

		tblCompanies = `
		CREATE TABLE IF NOT EXISTS companies (
			"id"           uuid PRIMARY KEY,
			"name"         text NOT NULL,
			"sync_addr"    text UNIQUE NOT NULL, /* primary server's sync_addr (but still want it unique) */
			"sync_api_key" text,                 /* used by secondary server */
			"sync_user_id" int,                  /* sync_user_fk constraint (can be NULL) */
			"active"       boolean,
			"created_at"   timestamp,
			"updated_at"   timestamp
		);`

		idxCompanies = `
		CREATE UNIQUE INDEX ON companies (lower(name));
		`
		tblSlices = `
		CREATE TABLE IF NOT EXISTS "slices" (
			"id"            uuid PRIMARY KEY,
			"name"          text NOT NULL,
			"slice_type"    slice_type_enum NOT NULL,
			"allow_sync"    boolean,
			"content_hash"  text,
			"content_count" integer,
			"content_date"  timestamp,
			"created_at"    timestamp,
			"updated_at"    timestamp
		);`

		idxSlices = `
		CREATE UNIQUE INDEX ON slices (lower(name));
    `
		tblSliceMetadata = `
		CREATE TABLE IF NOT EXISTS "slice_metadata" (
			"slice_id" uuid REFERENCES "slices" ON DELETE CASCADE,
			"key"      text,
			"value"    text,
			PRIMARY KEY ("slice_id", "key")
		);`

		tlbTags = `
		CREATE TABLE IF NOT EXISTS "tags" (
			"id"          serial PRIMARY KEY,
			"name"        text UNIQUE NOT NULL, /* lowercase by convention */
			"description" text,
			"created_at"  timestamp,
			"updated_at"  timestamp
		);`

		tblSliceTags = `
		CREATE TABLE IF NOT EXISTS "slice_tags" (
			"tag_id" int REFERENCES "tags" ON DELETE CASCADE,
			"slice_id" uuid REFERENCES "slices" ON DELETE CASCADE,
			PRIMARY KEY ("tag_id", "slice_id")
		);`

		tblSubscriptions = `
		CREATE TABLE IF NOT EXISTS "subscriptions" (
			"sub_id"       uuid PRIMARY KEY,
			"slice_id"     uuid REFERENCES "slices" ON DELETE CASCADE,
			"company_id"   uuid REFERENCES "companies" ON DELETE CASCADE,
			"name"         text NOT NULL,
			"description"  text,
			"active"       boolean,
			"created_at"   timestamp,
			"updated_at"   timestamp,
			CONSTRAINT "sub_alt_key" UNIQUE("slice_id", "company_id")
		);`

		idxSubscriptions = `
		CREATE UNIQUE INDEX ON subscriptions (lower(name));`

		tblGrains = `
		CREATE TABLE IF NOT EXISTS "grains" (
		  "id"           uuid PRIMARY KEY,
		  "slice_id"     uuid REFERENCES "slices" ON DELETE CASCADE,
		  "grain_key"    text NOT NULL,
		  "encoding"     encoding_enum,
		  "payload"      text,
		  "source"       text,
		  "created_at"   timestamp,
		  CONSTRAINT "grains_sliceid_grainkey_key" UNIQUE("slice_id", "grain_key")
		);`

		tblActivity = `
		CREATE TABLE IF NOT EXISTS "activity" (
		  "id"         serial PRIMARY KEY,
		  "sub_id"     uuid REFERENCES "subscriptions" ON DELETE CASCADE,
		  "success"    boolean,
		  "message"    text NOT NULL,
		  "duration"   timestamp,
		  "created_at" timestamp
		);`

		tblUsers = `
		CREATE TABLE IF NOT EXISTS users (
		  "id"               serial PRIMARY KEY,
		  "username"         text UNIQUE NOT NULL,
		  "password"         text,
		  "email"            text NOT NULL,
		  "first_name"       text NOT NULL,
		  "last_name"        text NOT NULL,
		  "phone"            text,
		  "active"           boolean,
		  "last_login"       timestamp,
		  "password_changed" timestamp,
		  "token"            text,
		  "role"             integer,
		  "company_id"       uuid REFERENCES "companies" ON DELETE RESTRICT,
		  "created_at"       timestamp,
		  "updated_at"       timestamp
		);`

		tblSettings = `
		CREATE TABLE IF NOT EXISTS "settings" (
		  "id" bool PRIMARY KEY DEFAULT TRUE, /* only allow one row */
			"server_role" server_role_enum,
			"server_id" uuid REFERENCES "companies" ON DELETE RESTRICT,
			"created_at"       timestamp,
			"updated_at"       timestamp,
			CONSTRAINT "settings_singleton" CHECK (id) /* only TRUE allowed */
			);`

		altCompanies = `
		ALTER TABLE companies
		ADD CONSTRAINT sync_user_fk FOREIGN KEY (sync_user_id) REFERENCES "users" ON DELETE RESTRICT;`
	)

	return []darwin.Migration{
		{Version: 1.00, Description: "Create Type '_enums'", Script: enums},
		{Version: 1.01, Description: "Create Table 'companies'", Script: tblCompanies},
		{Version: 1.02, Description: "Create Indexes on 'companies'", Script: idxCompanies},
		{Version: 1.03, Description: "Create Table 'slices'", Script: tblSlices},
		{Version: 1.04, Description: "Create Indexes on 'slices'", Script: idxSlices},
		{Version: 1.05, Description: "Create Table 'slice_metadata'", Script: tblSliceMetadata},
		{Version: 1.06, Description: "Create Table 'tags'", Script: tlbTags},
		{Version: 1.07, Description: "Create Table 'slice_tags'", Script: tblSliceTags},
		{Version: 1.08, Description: "Create Table 'subscriptions'", Script: tblSubscriptions},
		{Version: 1.09, Description: "Create Indexes on 'subscriptions'", Script: idxSubscriptions},
		{Version: 1.10, Description: "Create Table 'grains'", Script: tblGrains},
		{Version: 1.11, Description: "Create Table 'activity'", Script: tblActivity},
		{Version: 1.12, Description: "Create Table 'users'", Script: tblUsers},
		{Version: 1.13, Description: "Create Table 'settings'", Script: tblSettings},
		{Version: 1.14, Description: "Alter Table 'companies'", Script: altCompanies},
	}
}

// Migrate applies any outstanding schema versions and returns version status message
func Migrate(psn string) (string, error) {
	var v1, v2 float64

	db, err := sql.Open("postgres", psn)
	if err != nil {
		return "", err
	}

	driver := darwin.NewGenericDriver(db, darwin.PostgresDialect{})

	schema := defineSchema()
	infoChan := make(chan darwin.MigrationInfo, len(schema))
	d := darwin.New(driver, schema, infoChan)

	v1 = currentVersion(d)

	if err := d.Migrate(); err != nil {
		return "", err
	}
	close(infoChan)

	v2 = v1
	for info := range infoChan {
		if info.Status == darwin.Applied {
			v2 = info.Migration.Version
		}
	}

	return changes(v1, v2), nil
}

// currentVersion reads all records from migration table to get the latest version
func currentVersion(d darwin.Darwin) float64 {
	var v, ver float64

	if infoList, err := d.Info(); err == nil {
		// get latest version
		for _, info := range infoList {
			v = info.Migration.Version
			if v > ver && info.Status == darwin.Applied {
				ver = v
			}
		}
	}
	return ver
}

func changes(v1, v2 float64) string {
	if v1 != v2 {
		return fmt.Sprintf("DB Version: %g (migrated from %g to %g)", v2, v1, v2)
	}
	return fmt.Sprintf("DB Version: %g", v1)
}
