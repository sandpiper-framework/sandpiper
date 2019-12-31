/*
 * Database: sandpiper
 * Migration: Create initial database tables in an empty database
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
  "id"          uuid PRIMARY KEY,
  "name"        text NOT NULL,
  "active"      boolean,
  "created_at"  timestamp,
  "updated_at"  timestamp,
  "deleted_at"  timestamp
);

CREATE TABLE IF NOT EXISTS "slices" (
  "id"            uuid PRIMARY KEY,
  "name"          text,
  "content_hash"  text,
  "content_count" integer,
  "last_update"   timestamp,
  "created_at"    timestamp,
  "updated_at"    timestamp,
  "deleted_at"    timestamp
);

CREATE TABLE IF NOT EXISTS "slice_metadata" (
  "slice_id" uuid,
  "key"      text,
  "value"    text,
  PRIMARY KEY ("slice_id", "key")
);

CREATE TABLE IF NOT EXISTS "subscriptions" (
  "slice_id"     uuid REFERENCES "slices" ("id"),
  "company_id"   uuid REFERENCES "companies" ("id"),
  "name"         text,
  "description"  text,
  "active"       boolean,
  "created_at"   timestamp,
  "updated_at"   timestamp,
  PRIMARY KEY ("slice_id", "company_id")
);

CREATE TABLE IF NOT EXISTS "grains" (
  "id"           uuid PRIMARY KEY,
  "slice_id"     uuid,
  "grain_type"   smallint,
  "payload"      text,
  "created_at"   timestamp
); 

CREATE TABLE IF NOT EXISTS users (
  "id"                   serial PRIMARY KEY,
  "first_name"           text,
  "last_name"            text,
  "username"             text,
  "password"             text,
  "email"                text,
  "phone"                text,
  "active"               boolean,
  "last_login"           timestamp,
  "password_changed"     timestamp,
  "token"                text,
  "role"                 integer,
  "company_id"           uuid REFERENCES "companies" ("id"),
  "created_at"           timestamp,
  "updated_at"           timestamp
);

COMMIT;