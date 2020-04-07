/*
 * Project: sandpiper
 * Database: sandpiper
 * Migration: Create initial database tables for server in an empty database
 * Date: 2019-11-16
 */
 
BEGIN;

CREATE TYPE slice_type_enum AS ENUM (
  'aces-file',
  'aces-items',
  'asset-files',
  'pies-file',
  'pies-items',
  'pies-marketcopy',
  'pies-pricesheet',
  'partspro-file',
  'partspro-items'
);

CREATE TYPE encoding_enum AS ENUM (
  'raw',
  'b64',
  'z64'
  'a85'
  'z85'
);

CREATE TABLE IF NOT EXISTS "settings" (
  "id"    SERIAL PRIMARY KEY,
  "key"   text UNIQUE NOT NULL,
  "value" text NOT NULL
);

CREATE TABLE IF NOT EXISTS companies (
  "id"         uuid PRIMARY KEY,
  "name"       text NOT NULL,
  "sync_addr"  text,
  "active"     boolean,
  "created_at" timestamp,
  "updated_at" timestamp
);
CREATE UNIQUE INDEX ON companies (lower(name));

CREATE TABLE IF NOT EXISTS "slices" (
  "id"            uuid PRIMARY KEY,
  "name"          text NOT NULL,
  "slice_type"    slice_type_enum NOT NULL,
  "content_hash"  text,
  "content_count" integer,
  "content_date"  timestamp,
  "created_at"    timestamp,
  "updated_at"    timestamp
);
CREATE UNIQUE INDEX ON slices (lower(name));

CREATE TABLE IF NOT EXISTS "slice_metadata" (
  "slice_id" uuid REFERENCES "slices" ON DELETE CASCADE,
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
  "tag_id" int REFERENCES "tags" ON DELETE CASCADE,
  "slice_id" uuid REFERENCES "slices" ON DELETE CASCADE,
  PRIMARY KEY ("tag_id", "slice_id")
);

CREATE TABLE IF NOT EXISTS "subscriptions" (
  "sub_id"       serial PRIMARY KEY,
  "slice_id"     uuid REFERENCES "slices" ON DELETE CASCADE,
  "company_id"   uuid REFERENCES "companies" ON DELETE CASCADE,
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
  "slice_id"     uuid REFERENCES "slices" ON DELETE CASCADE,
  "grain_key"    text NOT NULL,
  "encoding"     encoding_enum,
  "payload"      text,
  "source"       text,
  "created_at"   timestamp,
  CONSTRAINT "grains_sliceid_grainkey_key" UNIQUE("slice_id", "grain_key")
);

CREATE TABLE IF NOT EXISTS "syncs" (
  "id"         serial PRIMARY KEY,
  "slice_id"   uuid REFERENCES "slices" ON DELETE CASCADE,
  "message"    text,
  "duration"   timestamp,
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
  "company_id"       uuid REFERENCES "companies",
  "created_at"       timestamp,
  "updated_at"       timestamp
);

COMMIT;