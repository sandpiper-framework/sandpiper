# Database Migrations

This directory contains SQL files for migrating the database from one version to the next.

The database version is checked each time the server is started and "migrated" to the latest version as required. If a
problem is encountered during a migration, the database will be marked "dirty", and you will need to correct the problem
manually. All migrations are tested before a new version is released, so we do not expect this to happen.

It should not be necessary, but a standalone migration tool (CLI) is available to force-apply corrected migrations
(and thus reset the "dirty" flag). Platform-specific CLI tools and usage instructions can be found here:

https://github.com/golang-migrate/migrate/tree/master/cmd/migrate

Download the correct executable (Linux or Windows) to the correct `cmd/migrations` directory and use `./migrations` as the
migration "source". For example:

```
# URI format: postgres://[user[:password]@][netloc][:port][/dbname][?param1=value1&...]

$ URI="postgres://sandpiper:autocare@localhost:5432/sandpiper?sslmode=disable"
$ migrate -source file://migrations $URI up 2
```

**Note that the "version" of a "dirty" database is the version that failed!** This is because it might have
applied some changes, but then encountered a problem that kept it from completing.

Note: `db_create.sql` is included in the migrations directory but is **not** applied by migrations (because it doesn't
have a number prefix). This sql script is used by Taskfile to create the initial database.