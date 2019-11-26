/*
 * Database: sandpiper
 * Migration: Create initial database tables in an empty database
 * Direction: Up
 * Date: 2019-11-16
 */
 
BEGIN; 

CREATE TABLE IF NOT EXISTS companies (
    id          serial PRIMARY KEY,
    name        varchar(30) NOT NULL,
    active      boolean,
    created_at  timestamp,
    updated_at  timestamp,
    deleted_at  timestamp
);

CREATE TABLE IF NOT EXISTS locations (
    id          serial PRIMARY KEY,
    name        text,
    address     text,
    active      boolean,
    company_id  integer REFERENCES companies (id),
    created_at  timestamp,
    updated_at  timestamp,
    deleted_at  timestamp
);

CREATE TABLE IF NOT EXISTS public.roles(
    id           serial PRIMARY KEY,
    access_level integer,
    name         text
);

CREATE TABLE IF NOT EXISTS users (
    id                   serial PRIMARY KEY,
    first_name           text,
    last_name            text,
    username             text,
    password             text,
    email                text,
    mobile               text,
    phone                text,
    address              text,
    active               boolean,
    last_login           timestamp,
    last_password_change timestamp,
    token                text,
    role_id              integer REFERENCES roles (id),
    company_id           integer REFERENCES companies (id),
    location_id          integer REFERENCES locations (id),
    created_at           timestamp,
    updated_at           timestamp,
    deleted_at           timestamp
);

COMMIT;