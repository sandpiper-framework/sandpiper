/*
 * Database: sandpiper
 * Migration: Create initial database tables in an empty database
 * Direction: Up
 * Date: 2019-11-16
 */
 
BEGIN;

CREATE TYPE "payloadtype" AS ENUM (
  'aces-file',
  'aces-item',
  'partspro-file',
  'partspro-item',
  'pies-file',
  'pies-item',
  'pies-marketcopy',
  'pies-pricesheet'
);

CREATE TABLE IF NOT EXISTS "settings" (
  "id"    SERIAL PRIMARY KEY,
  "key"   text UNIQUE NOT NULL,
  "value" text
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

CREATE TABLE IF NOT EXISTS "subscribers" (
  "id" uuid PRIMARY KEY,
  "company" text,
  "contact" text,
  "email"   text,
  "active"  boolean
);

CREATE TABLE IF NOT EXISTS "subscriptions" (
  "slice_id"      uuid REFERENCES "slices" ("id"),
  "subscriber_id" uuid REFERENCES "subscribers" ("id"),
  "name"          text,
  "description"   text,
  PRIMARY KEY ("slice_id", "subscriber_id")
);

CREATE TABLE IF NOT EXISTS "data_objects" (
  "id"           uuid PRIMARY KEY,
  "slice_id"     uuid,
  "payload_type" payloadtype,
  "payload"      text
); 

CREATE TABLE IF NOT EXISTS companies (
  "id"          serial PRIMARY KEY,
  "name"        varchar(30) NOT NULL,
  "active"      boolean,
  "created_at"  timestamp,
  "updated_at"  timestamp,
  "deleted_at"  timestamp
);

CREATE TABLE IF NOT EXISTS locations (
  "id"          serial PRIMARY KEY,
  "name"        text,
  "address"     text,
  "active"      boolean,
  "company_id"  integer REFERENCES "companies" ("id"),
  "created_at"  timestamp,
  "updated_at"  timestamp,
  "deleted_at"  timestamp
);

CREATE TABLE IF NOT EXISTS roles(
  "id"           serial PRIMARY KEY,
  "access_level" integer,
  "name"         text
);

CREATE TABLE IF NOT EXISTS users (
  "id"                   serial PRIMARY KEY,
  "first_name"           text,
  "last_name"            text,
  "username"             text,
  "password"             text,
  "email"                text,
  "mobile"               text,
  "phone"                text,
  "address"              text,
  "active"               boolean,
  "last_login"           timestamp,
  "last_password_change" timestamp,
  "token"                text,
  "role_id"              integer REFERENCES "roles" ("id"),
  "company_id"           integer REFERENCES "companies" ("id"),
  "location_id"          integer REFERENCES "locations" ("id"),
  "created_at"           timestamp,
  "updated_at"           timestamp,
  "deleted_at"           timestamp
);

COMMIT;