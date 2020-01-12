/*
 * Project: sandpiper
 * Database: sandpiper
 * Migration: Create initial database tables for server in an empty database
 * Direction: Up
 * Date: 2019-11-16
 */
 
BEGIN;

CREATE TABLE IF NOT EXISTS "settings" (
  "id"    SERIAL PRIMARY KEY,
  "key"   text UNIQUE NOT NULL,
  "value" text
);

CREATE TABLE IF NOT EXISTS companies (
  "id"         uuid PRIMARY KEY,
  "name"       text UNIQUE NOT NULL,
  "sync_addr"  text,
  "active"     boolean,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE IF NOT EXISTS "slices" (
  "id"            uuid PRIMARY KEY,
  "name"          text UNIQUE NOT NULL,
  "content_hash"  text,
  "content_count" integer,
  "content_date"  timestamp,
  "created_at"    timestamp,
  "updated_at"    timestamp
);

CREATE TABLE IF NOT EXISTS "slice_metadata" (
  "slice_id" uuid REFERENCES "slices" ON DELETE CASCADE,
  "key"      text,
  "value"    text,
  PRIMARY KEY ("slice_id", "key")
);

CREATE TABLE IF NOT EXISTS "subscriptions" (
  "slice_id"     uuid REFERENCES "slices" ON DELETE CASCADE,
  "company_id"   uuid REFERENCES "companies" ON DELETE CASCADE,
  "name"         text UNIQUE NOT NULL,
  "description"  text,
  "active"       boolean,
  "created_at"   timestamp,
  "updated_at"   timestamp,
  PRIMARY KEY ("slice_id", "company_id")
);

CREATE TABLE IF NOT EXISTS "grains" (
  "id"           uuid PRIMARY KEY,
  "slice_id"     uuid REFERENCES "slices" ON DELETE CASCADE,
  "grain_type"   smallint,
  "grain_key"    text,
  "payload"      bytea,
  "created_at"   timestamp,
  CONSTRAINT "grain_type_key" UNIQUE("slice_id", "grain_type", "grain_key")
); 

CREATE TABLE IF NOT EXISTS users (
  "id"               serial PRIMARY KEY,
  "username"         text UNIQUE NOT NULL,
  "password"         text,
  "email"            text,
  "first_name"       text,
  "last_name"        text,
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