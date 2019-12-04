# Database Migrations

This directory contains SQL files for migrating the database from one version to the next (and back if necessary).

The database version is checked each time the server is started and "migrated" to the latest version as required. If a
problem is encountered during a migration, the database will be marked "dirty", and you will need to correct the problem
manually.

A standalone migration tool (CLI) is available to move up/down by version steps. Platform-specific CLI tools and usage
instructions can be found here:

https://github.com/golang-migrate/migrate/tree/master/cmd/migrate

Download the correct executable (Linux or Windows) to the `cmd/migration` directory and use `./migrations` as the
migration "source". For example:

```
# URI format: postgres://[user[:password]@][netloc][:port][/dbname][?param1=value1&...]

$ URI="postgres://sandpiper:autocare@localhost:5432/sandpiper?sslmode=disable"
$ migrate -source file://migrations $URI up 2
```

**Note:** `db_create.sql` is included in the migrations directory but is **not** applied by migrations (because it doesn't
have a number prefix). This sql script is used by Taskfile to create the initial database.