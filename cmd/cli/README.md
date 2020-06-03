# Sandpiper CLI utility

## Introduction

Level 1 implementation of Sandpiper means syncing "file-based" data-objects. These `sandpiper` CLI tools support adding, pulling and listing Level 1 files.
This functionality is mostly duplicated by web-based Admin screens, but a command line interface allows an important alternative for hands-free automation.

### Command Syntax

```
sandpiper [global-options] [add | list | pull | sync | help] [command-options] <arguments>

global-options (also available via [environment variable]):
   --user name, -u name        server login name [SANDPIPER_USER]
   --password value, -p value  user password [SANDPIPER_PASSWORD]
   --config FILE, -c FILE      Load configuration from FILE (default: config.yaml) [SANDPIPER_CONFIG]
   --debug, -d                 provide some debugging information to stdout
   --help, -h                  show general help (default: false)
   --version, -v               print the version (default: false)

commands:
   add      add a file-based grain from a local file
   pull     save file-based grains to the file system
   list     list slices (if no slice provided) or file-based grains by slice_id or slice_name
   sync     start the sync process on active subscriptions
   init     initialize a sandpiper primary or secondary database
   secrets  generate  new random secrets for env vars and api-config.yaml file 
   help, h  Shows a list of commands or help for one command
```

If the password is not provided, you will be prompted for it.

### Configuration file

All commands look for a `cli-config.yml` file in the current directory to determine the sandpiper server address (this can be overridden by the `--config` option or by the SANDPIPER_CONFIG environment variable). 

Here is an example of the config.yml file.

```
command:
  url: http://localhost
  port: 8080
  max_sync_proceses: 5
```

## Initialize a Sandpiper database

The `init` command will create, migrate and seed a new primary or secondary postgresql database. (Please see the "Testing Workbook" for slightly different instructions if not installing for production).

When you issue this command, you will be prompted for your PostgreSQL Host Address, Port and Superuser credentials. This is required to create a new database. In most cases, you can simply press `Enter` to accept the default value (shown in parentheses).

```
sandpiper (v0.1.2-75-gd49f4eb)
Copyright 2020 The Sandpiper Authors. All rights reserved.

INITIALIZE A SANDPIPER DATABASE

PostgreSQL Address (localhost):
PostgreSQL Port (5432):
PostgreSQL Superuser (postgres):
PostgreSQL Superuser Password: *********
SSL Mode (disable):
connected to host
```
The `localhost` address (which is equivalent to 127.0.0.1) indicates you're running this command on the same machine as PostgreSQL. Otherwise, it would be a standard ip4 address on your network (e.g. 192.168.1.100) or possibly a hosted instance endpoint (e.g. myinstance.123456789012.us-east-1.rds.amazonaws.com). The superuser password (from above) will be hidden when you type.

You should see "connected to host" to indicate that the connection was successful. Next, you will be prompted for the new database information.

```
New Database Name (sandpiper):
Database Owner (sandpiper):
Database Owner Password: focal-weedy-brood-hat
CREATE DATABASE sandpiper;
CREATE USER sandpiper WITH ENCRYPTED PASSWORD 'focal-weedy-brood-hat';
GRANT ALL PRIVILEGES ON DATABASE sandpiper TO sandpiper;

applying migrations...
Database: "sandpiper"
DB Version: 1.15 (migrated from 0.00 to 1.15)
```

The recommended database name is `sandpiper` regardless of the server-role you require (primary or secondary). If you need both server-roles, you can name it anything you like ("secondary", "receiver", "tidepool", etc.).

In the example above, we used default values except when required to enter a password for the database owner (please select a strong password of your own) and again, keep a record of it for later. You will need it when starting the server.

The database owner is the only user to connect directly to the database (via the sandpiper-api server). This should not be confused with a sandpiper end-user which is stored in the `users` table for authentication and access.

```
Company Name: Better Brakes
Server-Role (primary*/secondary): primary
Public Sync URL: http://localhost:8080
Server http URL (http://localhost): 
Added Company "Better Brakes"

Sandpiper Admin Password: admin
Added User "admin"

initialization complete for "sandpiper"

Server config file "api-primary.yaml" created in C:\sandpiper
Command config file "cli-primary.yaml" created in C:\sandpiper
```

In production, you would enter a strong admin password, but enter "admin" here to make testing easier. Also, the public sync URL would normally be something like `https://sandpiper.betterbrakes.com`, but we are going to test locally, on the same machine with both servers (using two different "ports").

## Add File-Based Objects

The `add` command creates a "file" data-object (i.e. grain) and adds it to a slice. This command could be called by an internal PIM, for example, to "publish" completed delivery files. By convention,
all L1 grains have a grain_key of "level-1" and use "z64" encoding (gzip/base64).

### Usage

```
sandpiper [global-options] add [command-options] <unzipped-file-to-add>

command-options:
   --slice value, -s value  either a slice_id (uuid) or slice_name (case-insensitive)
   --noprompt               do not prompt before over-writing a grain (default is to prompt)

arguments:
    A single filename (absolute or relative to the command) that should be added to the provided slice.
    The file should in its native format (e.g. .xml, .txt) and should *not* be zipped. Compression is
    performed separately by the add command.

Example:
    sandpiper -u user -p password add --slicename "aap-brake-pads" --noprompt acme_brakes_full_2019-12-12.xml
    sandpiper -u admin -p admin add --slice 2bea8308-1840-4802-ad38-72b53e31594c testdata\aces-file.xml
```

This command adds the ACES xml file as a grain as defined by the supplied request body (see below).
  
### Sandpiper API

The following api is used to add an object with the sandpiper add command:

```
POST /v1/grains
```

The request "body" submitted by `sandpiper add` would look like this:

```
{
    "slice_id": "2bea8308-1840-4802-ad38-72b53e31594c",
    "grain_key": "level-1",
    "encoding": "z85"
    "payload": "--gzip/ascii85 encoded data from supplied filename--"
}
```

The response would include the assigned grain "id" (but not the payload). Use the `sandpiper list --full` command to display all information stored.

```
{
  "id": "cb5872d2-6f98-4f23-bd2f-7366b8063523",
  "slice_id": "1b40204a-7acd-4c78-a3c4-0fa95d2f00f6",
  "grain_key": "level-1",
  "source": "acme_brakes_full_2019-12-12.xml",
  "encoding": "z85",
  "payload_len": 2696,
  "created_at": "2020-04-07T15:59:40.5223569-05:00"
}
```

## List File-Based Objects

This command displays slice information or grain information for a slice.

#### Syntax:

```
sandpiper [global-options] list [command-options]

command-options:
   --slice value, -s value  either a slice_id (uuid) or slice_name (case-insensitive)
   --full                   provide full listings (default: false)
   --help, -h               show help (default: false)

    If a slice is not provided, a listing of slices is displayed to stdout.
    If a slice is provided and a valid UUID, it is interpreted as a slice_id.
    If a slice is provided and not a valid UUID, it is interpreted as a slice_name. 

Examples:
    sandpiper -u user -p password list
    sandpiper -u user -p password list --slice 1b40204a-7acd-4c78-a3c4-0fa95d2f00f6
    sandpiper -u user -p password list --full --slice aap-brake-pads
```

## Pull File-Based Objects

Implement the "pull" command to retrieve "file" data-objects from an optional slice in the pool. If the slice is not supplied it will create a sub-directory for each one it finds.

#### Syntax:

```
sandpiper [global-options] pull [command-options] <output-directory>

command-options:
   --slice value, -s value  either a slice_id (uuid) or slice_name (case-insensitive)
   --help, -h               show help (default: false)

    If a slice is not provided, all level-1 grains are extracted to the local file system.
    If a slice is provided and a valid UUID, it is interpreted as a slice_id.
    If a slice is provided and not a valid UUID, it is interpreted as a slice_name.  

arguments:
   optional output directory (will be created if it doesn't exist)

Examples:
    sandpiper -u user -p password pull --slice 1b40204a-7acd-4c78-a3c4-0fa95d2f00f6
    sandpiper -u user -p password pull --slice "aap-brake-pads" /temp/output
    sandpiper -u uers -p password pull /temp/output

    This last example will create a directory structure under /temp/output with the slice name as the parent
    for each grain pulled:

    /temp/output
        |-- slice1
            |-- grain1
        |-- slice2
            |-- grain2 
```
    
Use forward slashes for the output directory even if on Windows. If an an output directory is not supplied as an argument, the current directory will be used.

The grain `source` field is used as the filename when saving the payload. If this value is empty, the slice-id is used instead.

## Sync Our Subscriptions

The sync command is run by an admin from a secondary server. It connects to each company with a sync_addr and retrieves our subscriptions. If a new subscription
is found we add it locally. If a subscription is disabled on the Primary, disable it locally and log the activity. If enabled on the Primary but not on our server,
we do not make any changes. We then perform a grain sync on all unlocked active slices assigned to that subscription.  

#### Syntax:

```
sandpiper [global-options] sync [command-options]

command-options:
   --partner value, -p value  limit to company name (case-insensitive) or company_id
   --noupdate                 Perform the sync without actually changing anything locally (default: false)
   --help, -h                 show help (default: false)
```

## Generate API Secrets

```
sandpiper secrets
```

This command will generate random values suitable for API secrets. Secrets are generated automatically by `sandpiper init`, but could be useful if you want to reset your secrets later.
