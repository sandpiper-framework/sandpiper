# Database Migrations

This folder contains SQL files for migrating the database from one version to the next (and back if necessary).

The database version is checked each time the server is started and migrated as required. If a problem is encountered with
a migration, you may want to use the standalone migration tool (CLI) to move up/down by version.

Platform specific CLI tools and usage instructions can be found here:

https://github.com/golang-migrate/migrate/tree/master/cmd/migrate