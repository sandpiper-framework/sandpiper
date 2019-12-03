# Database Migrations

This folder contains SQL files for migrating the database from one version to the next (and back if necessary).

The database version is checked each time the server is started and "migrated" to the latest version as required. If a
problem is encountered during a migration, the database will be marked "dirty", and you will need to correct the problem
manually. A standalone migration tool (CLI) is available to move up/down by version steps.

Platform specific CLI tools and usage instructions can be found here:

https://github.com/golang-migrate/migrate/tree/master/cmd/migrate

**Note:** `db_create.sql` is included in the migrations directory but is **not** applied by migrations (because it doesn't
have a number prefix). This sql script is used by Taskfile to create the initial database.