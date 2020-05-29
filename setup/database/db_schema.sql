/*
 Copyright The Sandpiper Authors. All rights reserved. Use of this source code is governed by an
 MIT-style license that can be found in the LICENSE.md file.

 This script is provided for documentation purposes only. See the README.md file for more information.

 Date: 2020-05-18
 DB Version 1.14
*/
 
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

CREATE TABLE IF NOT EXISTS companies (
  "id"           uuid PRIMARY KEY,
  "name"         text NOT NULL,
  "sync_addr"    text UNIQUE NOT NULL, /* primary server's sync_addr (but still want it unique) */
  "sync_api_key" text,                 /* used by secondary server */
  "sync_user_id" int,                  /* sync_user_fk constraint (can be NULL) */
  "active"       boolean,
  "created_at"   timestamp,
  "updated_at"   timestamp
);
CREATE UNIQUE INDEX ON companies (lower(name));

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
);
CREATE UNIQUE INDEX ON slices (lower(name));

CREATE TABLE IF NOT EXISTS "slice_metadata" (
  "slice_id" uuid REFERENCES "slices" ("id") ON DELETE CASCADE,
  "key"      text,
  "value"    text,
  PRIMARY KEY ("slice_id", "key")
);

CREATE TABLE IF NOT EXISTS "tags" (
  "id"          serial PRIMARY KEY,
  "name"        text UNIQUE NOT NULL, /* lowercase by convention */
  "description" text,
  "created_at"  timestamp,
  "updated_at"  timestamp
);

CREATE TABLE IF NOT EXISTS "slice_tags" (
  "tag_id" int REFERENCES "tags" ("id") ON DELETE CASCADE,
  "slice_id" uuid REFERENCES "slices" ("id") ON DELETE CASCADE,
  PRIMARY KEY ("tag_id", "slice_id")
);

CREATE TABLE IF NOT EXISTS "subscriptions" (
  "sub_id"       uuid PRIMARY KEY,
  "slice_id"     uuid REFERENCES "slices" ("id") ON DELETE CASCADE,
  "company_id"   uuid REFERENCES "companies" ("id") ON DELETE CASCADE,
  "name"         text NOT NULL,
  "description"  text,
  "active"       boolean,
  "created_at"   timestamp,
  "updated_at"   timestamp,
  CONSTRAINT "sub_alt_key" UNIQUE("slice_id", "company_id")
);
CREATE UNIQUE INDEX ON subscriptions (lower(name));

CREATE TABLE IF NOT EXISTS "grains" (
  "id"           uuid PRIMARY KEY,
  "slice_id"     uuid REFERENCES "slices" ("id") ON DELETE CASCADE,
  "grain_key"    text NOT NULL,
  "encoding"     encoding_enum,
  "payload"      text,
  "source"       text,
  "created_at"   timestamp,
  CONSTRAINT "grains_sliceid_grainkey_key" UNIQUE("slice_id", "grain_key")
);

CREATE TABLE IF NOT EXISTS "activity" (
  "id"         serial PRIMARY KEY,
  "sub_id"     uuid REFERENCES "subscriptions" ("sub_id") ON DELETE CASCADE,
  "success"    boolean,
  "message"    text NOT NULL,
  "error"      text,
  "duration"   bigint,
  "created_at" timestamp
);

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
  "company_id"       uuid REFERENCES "companies" ("id") ON DELETE RESTRICT,
  "created_at"       timestamp,
  "updated_at"       timestamp
);

CREATE TABLE IF NOT EXISTS "settings" (
  "id" bool PRIMARY KEY DEFAULT TRUE, /* only allow one row */
  "server_role" server_role_enum,
  "server_id" uuid REFERENCES "companies" ("id") ON DELETE RESTRICT,
  "created_at"       timestamp,
  "updated_at"       timestamp,
  CONSTRAINT "settings_singleton" CHECK (id) /* only TRUE allowed */
);

ALTER TABLE companies
    ADD CONSTRAINT sync_user_fk FOREIGN KEY (sync_user_id) REFERENCES "users" ("id") ON DELETE RESTRICT;
