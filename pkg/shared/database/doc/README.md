# Creating Primary or Secondary Database

The `sandpiper init` command should be used to create a new primary or secondary database. This process will also initialize the database (admin user, etc.) and create a config.yaml file for the server to use.

# Manually Creating a Primary or Secondary Database

Change the dbname, dbuser and variables in the `db_create.sql` file and remove either database (primary or secondary) that is not required in your
environment. Run `task db-build` to actually create the database.

NOTE: You will still need to run `sandpiper init` to migrate the database and seed it with required records.

# Database Migrations

The database version is checked each time the server is started and "migrated" to the latest version as required.
