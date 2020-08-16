# Database Setup

This directory contains information about the Sandpiper database schema and migration process. 

Sandpiper requires a PostgreSQL sever. Please see the appropriate setup document for help getting started with PostgreSQL.

# Creating Primary or Secondary Database

The `sandpiper init` command should be used to create a new primary or secondary database, build the tables and initialize the database (with admin user, etc.). It will also create configuration files for the server and sandpiper command line program to use.

# Database Migrations

The database version is checked each time the server is started and "migrated" to the latest version as required. These migrations are stored in a source code file (schema.go) and so are not available as individual sql files.
Each migration "group" is given a major version number (1.xx, 2.xx) with minor numbers (x.01, x.02) representing an actual migration step. When a migration is completed, rows are added to the "darwin_migrations" table indicating that they were successfully applied. 

# Manually Creating a Primary or Secondary Database

Even though the controlling schema is contained in the source code, a full DDL script of the database is provided (db_schema.sql) for documentation purposes. This file is updated with each release and is intended to represent the current state (as shown by the db version number in the header). 

If for some reason you need to create the database manually, change the dbname, dbuser and variables in the `db_create.sql` file and remove either database (primary or secondary) that is not required in your
environment. Run `task db-build` to create the database(s).

NOTE: You will still need to run `sandpiper init` to migrate the database (as explained above) and "seed" it with required records.

