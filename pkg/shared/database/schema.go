// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package database

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/GuiaBolso/darwin"
	"github.com/lib/pq"
)

// defineSchema returns a list of our database migrations
// Each migration script is defined in a separate descriptive string variable (versioned by major db release)
// prefixes are "tbl" (create table), "idx" (create index), "alt" (alter)
// comments must only appear after sql on a line and cannot span lines (comments are stripped before checksum calc)
// *NEVER* change/remove a step once released! (because a checksum of the script is saved with the migration)
func defineSchema() []darwin.Migration {
	var (
		enumsV1 = `
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
		);
		CREATE TYPE sync_status_enum AS ENUM (
			'none',
			'updating',
			'success',
			'error'
		);`

		tblCompaniesV1 = `
		CREATE TABLE IF NOT EXISTS companies (
			"id"           uuid PRIMARY KEY,
			"name"         text NOT NULL,
			"sync_addr"    text UNIQUE NOT NULL, /* primary server's sync_addr (but still want it unique) */
			"sync_api_key" text,                 /* used by secondary server */
			"sync_user_id" int,                  /* only on primary (can be NULL) sync_user_fk constraint */
			"active"       boolean,
			"created_at"   timestamp,
			"updated_at"   timestamp
		);`

		idxCompaniesV1 = `
		CREATE UNIQUE INDEX ON companies (lower(name));`

		tblSlicesV1 = `
		CREATE TABLE IF NOT EXISTS "slices" (
			"id"                uuid PRIMARY KEY,
			"name"              text NOT NULL,
			"slice_type"        slice_type_enum NOT NULL,
			"content_count"     integer,
			"content_date"      timestamp,
			"content_hash"      text,
			"allow_sync"        boolean,                    /* locked during content update */
			"sync_status"       sync_status_enum NOT NULL,  /* only on secondary */
			"last_sync_attempt" timestamp,                  /* only on secondary */
			"last_good_sync"    timestamp,                  /* only on secondary */
			"created_at"        timestamp,
			"updated_at"        timestamp
		);`

		idxSlicesV1 = `
		CREATE UNIQUE INDEX ON slices (lower(name));`

		tblSliceMetadataV1 = `
		CREATE TABLE IF NOT EXISTS "slice_metadata" (
			"slice_id" uuid REFERENCES "slices" ON DELETE CASCADE,
			"key"      text,
			"value"    text,
			PRIMARY KEY ("slice_id", "key")
		);`

		tlbTagsV1 = `
		CREATE TABLE IF NOT EXISTS "tags" (
			"id"          serial PRIMARY KEY,
			"name"        text UNIQUE NOT NULL, /* lowercase by convention */
			"description" text,
			"created_at"  timestamp,
			"updated_at"  timestamp
		);`

		tblSliceTagsV1 = `
		CREATE TABLE IF NOT EXISTS "slice_tags" (
			"tag_id" int REFERENCES "tags" ON DELETE CASCADE,
			"slice_id" uuid REFERENCES "slices" ON DELETE CASCADE,
			PRIMARY KEY ("tag_id", "slice_id")
		);`

		tblSubscriptionsV1 = `
		CREATE TABLE IF NOT EXISTS "subscriptions" (
			"sub_id"       uuid PRIMARY KEY,
			"slice_id"     uuid REFERENCES "slices" ON DELETE RESTRICT,
			"company_id"   uuid REFERENCES "companies" ON DELETE RESTRICT,
			"name"         text NOT NULL,
			"description"  text,
			"active"       boolean,
			"created_at"   timestamp,
			"updated_at"   timestamp,
			CONSTRAINT "sub_alt_key" UNIQUE("slice_id", "company_id")
		);`

		idxSubscriptionsV1 = `
		CREATE UNIQUE INDEX ON subscriptions (lower(name));`

		tblGrainsV1 = `
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

		tblActivityV1 = `
		CREATE TABLE IF NOT EXISTS "activity" (
			"id"         serial PRIMARY KEY,
			"company_id" uuid REFERENCES "companies" ON DELETE CASCADE,
			"sub_id"     uuid REFERENCES "subscriptions" ON DELETE CASCADE,
			"success"    boolean,
			"message"    text NOT NULL,
			"error"      text,
			"duration"   bigint,
			"created_at" timestamp
		);`

		tblUsersV1 = `
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

		tblSettingsV1 = `
		CREATE TABLE IF NOT EXISTS "settings" (
			"id" bool PRIMARY KEY DEFAULT TRUE, /* only allow one row */
			"server_role" server_role_enum,
			"server_id" uuid REFERENCES "companies" ON DELETE RESTRICT,
			"created_at"       timestamp,
			"updated_at"       timestamp,
			CONSTRAINT "settings_singleton" CHECK (id) /* only TRUE allowed */
		);`

		altCompaniesV1 = `
		ALTER TABLE companies
		ADD CONSTRAINT sync_user_fk FOREIGN KEY (sync_user_id) REFERENCES "users" ON DELETE RESTRICT;`
	) // v1 release
	var (
	// this is a placeholder to show our change pattern of one release per var(...) making code-folding easier.
	) // v2 release

	// minify simplifies the script to keep certain changes (spaces, tabs, case and comments) from creating a new checksum
	var minify = func(script string) string {
		b := strings.Builder{}
		s := strings.ToLower(strings.ReplaceAll(script, "/*", "--"))
		lines := strings.Split(s, "\n")
		for _, line := range lines {
			if i := strings.Index(line, "--"); i != -1 {
				line = line[0:i]
			}
			b.WriteString(strings.TrimSpace(line) + "\n")
		}
		result := strings.TrimSpace(strings.ReplaceAll(b.String(), "\t", " "))
		before := 0
		for len(result) != before {
			before = len(result)
			result = strings.ReplaceAll(result, "  ", " ")
		}
		return strings.TrimSpace(result)
	}

	// Each database change release is given a major version number (1.xx, 2.xx) with minor numbers (x.01, x.02)
	// representing the actual migration steps within that release. Version numbers must be ascending with a
	// convention to skip x.00 (i.e. "steps" start from 01).
	return []darwin.Migration{
		{Version: 1.01, Description: "Create Type '_enums'", Script: minify(enumsV1)},
		{Version: 1.02, Description: "Create Table 'companies'", Script: minify(tblCompaniesV1)},
		{Version: 1.03, Description: "Create Indexes on 'companies'", Script: minify(idxCompaniesV1)},
		{Version: 1.04, Description: "Create Table 'slices'", Script: minify(tblSlicesV1)},
		{Version: 1.05, Description: "Create Indexes on 'slices'", Script: minify(idxSlicesV1)},
		{Version: 1.06, Description: "Create Table 'slice_metadata'", Script: minify(tblSliceMetadataV1)},
		{Version: 1.07, Description: "Create Table 'tags'", Script: minify(tlbTagsV1)},
		{Version: 1.08, Description: "Create Table 'slice_tags'", Script: minify(tblSliceTagsV1)},
		{Version: 1.09, Description: "Create Table 'subscriptions'", Script: minify(tblSubscriptionsV1)},
		{Version: 1.10, Description: "Create Indexes on 'subscriptions'", Script: minify(idxSubscriptionsV1)},
		{Version: 1.11, Description: "Create Table 'grains'", Script: minify(tblGrainsV1)},
		{Version: 1.12, Description: "Create Table 'activity'", Script: minify(tblActivityV1)},
		{Version: 1.13, Description: "Create Table 'users'", Script: minify(tblUsersV1)},
		{Version: 1.14, Description: "Create Table 'settings'", Script: minify(tblSettingsV1)},
		{Version: 1.15, Description: "Add Foreign Key 'sync_user_fk'", Script: minify(altCompaniesV1)},
	}
}

// Migrate applies any outstanding schema versions and returns a version status message
func Migrate(dsn string) (string, error) {
	var v1, v2 float64

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return "", err
	}

	count, v1, err := currentVersion(db)
	if err != nil {
		return "", err
	}

	schema := defineSchema()
	if count == len(schema) && v1 == schema[count-1].Version {
		// already up-to-date
		return changes(v1, v1), nil
	}

	// setup for the migrations
	driver := darwin.NewGenericDriver(db, darwin.PostgresDialect{})
	infoChan := make(chan darwin.MigrationInfo, len(schema)) // for error reporting
	d := darwin.New(driver, schema, infoChan)

	// perform the migrations
	if err := d.Migrate(); err != nil {
		close(infoChan)
		_, v2, _ = currentVersion(db)
		return "", fmt.Errorf("migration (was v%.2f now v%.2f): %w\n\n%s\n", v1, v2, err, progress(infoChan))
	}
	close(infoChan)

	_, v2, err = currentVersion(db)
	if err != nil {
		return "", err
	}

	return changes(v1, v2), nil
}

// currentVersion reads from migration table to get the latest version and number of steps applied
func currentVersion(db *sql.DB) (count int, ver float64, err error) {
	s := `select count(*) as n, max(version) as ver from darwin_migrations;`
	err = db.QueryRow(s).Scan(&count, &ver)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "42P01" {
				// missing table is not an error during first migration
				count, ver, err = 0, 0, nil
			}
		}
	}
	return count, ver, err
}

// progress returns the steps attempted during this migration
func progress(ch <-chan darwin.MigrationInfo) string {
	var b strings.Builder

	for info := range ch {
		_, _ = fmt.Fprintf(&b, "v%.2f: \"%s\" (%s) Error: %v\n",
			info.Migration.Version, info.Migration.Description, info.Status.String(), info.Error)
	}
	return b.String()
}

func changes(v1, v2 float64) string {
	if v1 != v2 {
		return fmt.Sprintf("DB Version: %.2f (migrated from %.2f to %.2f)", v2, v1, v2)
	}
	return fmt.Sprintf("DB Version: %g", v1)
}
