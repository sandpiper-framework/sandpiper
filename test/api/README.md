# Testing Handbook

This handbook walks you through the testing process. The goal is to set up two servers (publisher and receiver), add a subscription and data file to the publisher and have it sync up with the receiver.

Note that all directory paths are shown in Linux format (forward slash) since they also work in Windows PowerShell. Also, paths are shown relative to the installation directory.

## Download Sandpiper Distribution Package

https://github.com/sandpiper-framework/sandpiper/releases

Under "Assets", select the proper binaries for your platform (Windows or Linux) and unzip the contents to a local directory.

## Download and Install PostgreSQL

https://www.postgresql.org/download/

See our separate guides for installing PostgreSQL on specific platforms.

Be sure to write down the superuser (usually `postgres`) and password. You will need them to manage your PostgreSQL server.

## Creating the Primary Database

Before we can do anything, we need to create a sandpiper database within the PostgreSQL server. A simple command line tool is provided to take care of this for you. Open a command prompt (terminal) window and enter the following commands (assuming you are currently in the same folder as the `sandpiper` CLI).

```
./sandpiper init --id 10000000-0000-0000-0000-000000000000
```
Notice that we included the `--id` option on the `init` command. This option lets you provide the `server_id` rather than having the software assign a random unique value (allowing our existing tests to work without change). This should not be done in a production environment because we want each company_id to be globally unique.

You will be prompted for your PostgreSQL Host Address, Port and Superuser credentials. This is required to create a new database. In most cases, you can simply press `Enter` to accept the default value (shown in parentheses).

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

The recommended database name is `sandpiper` regardless of the server-role you require (primary or secondary). If you need both server-roles, you can name the second database anything you like ("secondary", "receiver", "tidepool", etc.).

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

## Creating the Secondary Database

Use the same procedure as above, but use this command to set the receiver's server-id:
 
 ```
 ./sandpiper init --id 20000000-0000-0000-0000-000000000000
 ```
 
Complete the PostgreSQL prompts as before. When prompted for the New Database Name, enter `tidepool` and accept the default Database Owner (sandpiper). Be sure to **enter the same sandpiper password you used above** (you will see a message that user "sandpiper" already exists, but that is not a problem as long as you use the same password).

You will then be prompted for the company information for the secondary database. Be sure to enter "secondary" for the Server-Role. It will not prompt for a "Public Sync URL" because its role is as a receiver":

```
Company Name: eCatCompany
Server-Role (primary*/secondary): secondary
Server http URL (http://localhost):
Added Company "eCatCompany"

Sandpiper Admin Password: admin
Added User "admin"

initialization complete for "tidepool"

Server config file "api-secondary.yaml" created in C:\sandpiper
Command config file "cli-secondary.yaml" created in C:\sandpiper
```
You should now have two databases, each with an "admin" user and an associated "company".

## Moving the Config files

The `sandpiper init` command creates two configuration files as its final step, one for the server ("api") and one for the `sandpiper` command line interface ("cli"). The files are named with the pattern (api,cli)-(server-role).yaml. So, for example, "api-primary.yaml" and "cli-primary.yaml".

Most installations will only use a single server role, but in our case we're running both on the same machine for testing purposes. When the API server starts, it looks for a file named "api-config.yaml" in the current directory, so we'd normally rename ours to match this default.

The `sandpiper` command looks for a file named `cli-config.yaml` in the current directory. Again, in our case, since we created two server-roles we have two separate API configuration files.

Next we'll run the sandpiper server (using the "primary" database) and create `subscriptions` and `grains` for us to sync. We'll do most of this work with a free REST client called Insomnia (someone must have thought that name was clever).

## Insomnia REST Client ("Core")

https://insomnia.rest/

- Insomnia is a cross-platform REST client for debugging APIs.
- It makes authentication easy with "chained" requests (able to extract values from the responses of other requests).
- You can save/load (checkin) and share setups easily
- Generate documentation from the insomnia file (https://github.com/jozsefsallai/insomnia-documenter)

After installing the software, [Import](https://support.insomnia.rest/article/52-importing-and-exporting-data) the workspace we provide called `Insomnia.json`. You should now see the "Sandpiper" workspace and several requests (organized in folders by resource).

Notice that we have two "environments" (Primary and Secondary). This allows us to change the "base_url" of the server we're accessing just by selecting the proper environment. Here, for example, we're sending all requests to `localhost` port 8080:

```
{
  "base_url": "localhost:8080",
  "jwt": "{% response 'body', 'req_0b8deb7ec8fd428a85b24f3e9a65300a', 'b64::JC50b2tlbg==::46b', 'never' %}"
}
```

The "jwt" is a variable we set up to allow "chained" requests. This is a very helpful feature because we need to use the JWT returned from a `/login` request on all subsequent requests.

## Starting the Primary Server

Before we can issue any requests to a server, however, we need to get one listening. This is very simple to do.

```
./api -config api-primary.yaml
```
If everything is working properly, you should see something like the following:

```
Sandpiper API Server (v0.1.2-75-gd49f4eb)
Copyright 2019-2020 The Sandpiper Authors. All rights reserved.

Database: "sandpiper"
DB Version: 1.15
Server role: "primary"
Server ID: 10000000-0000-0000-0000-000000000000

⇨ http server started on [::]:8080
```
That last line shows that a web server is running and listening on [http://localhost:8080/check](http://localhost:8080/check). You should be able to open a browser and type that address and receive a response of:

```
"Sandpiper API OK"
```

We also provide this "Health Check" as the first request in the Insomnia workspace. Select that request now and Press `Send`. You should see that it returns that body response along with a 200 http status code and 13 header rows.

## Authentication

Let's move on to adding some Companies and Users. Before we can do that, though, we need to issue a Login request to get the JWT authentication token.

Click on the `POST Login` request and notice that the following JSON is pre-configured in the body of that request:

```json
{ "username": "admin", "password":"admin" }
```
Assuming you didn't change anything when running the `sandpiper init` command, these credentials should work without modification. If you did change something, simply make those same changes here. Press the `Send` button to send the Login request to our sever. You should see something like this:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjIjoiNWMzN2QyZmQtY2EzZC00ZTVlLThkOGEtNDUwNmFbWlu...",
  "expires": "2020-05-18T16:27:43-05:00",
  "refresh_token": "d16e1a83587fd1546a204725c97dcca29c2c7422"
}
```

That token has now been saved in our Insomnia `jwt` variable, and we can use it to make authenticated requests (as admin). The user information is actually encrypted in the token itself so that the server knows who is making those requests.

If you look at the console where you're running the server, you will also notice that there is lots of information displayed whenever a request is made. There's a setting in the config file to supress some of these messages, but for now it can be helpful while testing to keep it.

## Companies and Users

Let's move on to adding some Companies and Users. Select the `List Companies` request and press `Send`. You should see something like the following:

```json
{
  "companies": [
    {
      "id": "10000000-0000-0000-0000-000000000000",
      "name": "Better Brakes",
      "sync_addr": "http://localhost:8080",
      "active": true,
      "created_at": "2020-05-18T22:45:52.273189Z",
      "updated_at": "2020-05-18T22:45:52.273189Z",
      "users": [
        {
          "id": 1,
          "first_name": "Sandpiper",
          "last_name": "Admin",
          "username": "admin",
          "email": "admin@mail.com",
          "active": true,
          "last_login": "2020-05-19T00:57:18.274078Z",
          "password_changed": "0001-01-01T00:00:00Z",
          "role": 100,
          "company_id": "10000000-0000-0000-0000-000000000000",
          "created_at": "2020-05-18T22:46:27.717492Z",
          "updated_at": "2020-05-19T00:57:18.275079Z"
        }
      ]
    }
  ],
  "paging": {
    "page_number": 1,
    "items_limit": 100,
    "items_total": 1
  }
}
```

This is the primary server's owner company. It currently has just one user (which is also the one you logged in under).

It's interesting to note that the company `id` is a globally unique identifier (because they can be shared across servers), but the user `id` is simply a sequential number (because they are local to this server).

A good example of a shared company is "eCatCompany" (which we created as the owner of our secondary server above). When a primary company agrees to share its product data, it adds a "secondary" company to its list of companies it does business with.
 
Let's do that now by selecting the `Add Copmany (eCat)` request and pressing `Send`. Then, let's add a user for our eCat company (`Add User (eCat)` under the User Resource). If you now go back and send the List Companies request again, you should see this second one added to the list:

```json
{
  "id": "20000000-0000-0000-0000-000000000000",
  "name": "eCat Co",
  "sync_addr": "http://localhost:8081",
  "active": true,
  "created_at": "2020-05-19T01:29:06.670001Z",
  "updated_at": "2020-05-19T01:29:06.670001Z",
  "users": [
    {
      "id": 2,
      "first_name": "Mickey",
      "last_name": "Mouse",
      "username": "companyadmin",
      "email": "mickey@gmail.com",
      "active": true,
      "last_login": "0001-01-01T00:00:00Z",
      "password_changed": "0001-01-01T00:00:00Z",
      "role": 120,
      "company_id": "20000000-0000-0000-0000-000000000000",
      "created_at": "2020-05-19T01:30:45.456748Z",
      "updated_at": "2020-05-19T01:30:45.456748Z"
    }
  ]
}
```

## Slices and Subscriptions

Now that we have two trading partners defined, we need a structure that organizes the data we want to share between them. This structure is called a `slice` and the way we assign these slices to companies is with a `subscription`.

Under the "Slice Resource" folder you should see two "POST Add" requests ("Add Slice1" & "Add Slice2"). In each case, the request body is provided to create a new slice (with metadata included in Slice2). Send "Add Slice2". You should see something like the following:

```json
{
  "id": "2bea8308-1840-4802-ad38-72b53e31594c",
  "slice_name": "Slice2",
  "slice_type": "aces-file",
  "content_hash": "",
  "content_count": 0,
  "content_date": "0001-01-01T00:00:00Z",
  "allow_sync": false,
  "sync_status": "none",
  "last_sync_attempt": "0001-01-01T00:00:00Z",
  "last_good_sync": "0001-01-01T00:00:00Z",
  "created_at": "2020-05-30T16:02:48.9188445-05:00",
  "updated_at": "2020-05-30T16:02:48.9188445-05:00",
  "metadata": {
    "pcdb-version": "2019-09-27",
    "vcdb-version": "2019-09-27"
  }
}
```
Adding a subscription assigns a slice to a company (note that a slice could actually be assigned to more than one company). Go ahead and find the "Add Subscription 1" request in the "Subscription Resource" folder and "Send" it. You should see a response like this:

```json
{
  "id": "2276e31d-7014-4947-b4ec-1bc170627593",
  "slice_id": "2bea8308-1840-4802-ad38-72b53e31594c",
  "company_id": "20000000-0000-0000-0000-000000000000",
  "name": "Subscription 1 (company2, slice2)",
  "description": "",
  "active": true,
  "created_at": "2020-05-30T16:17:34.8817721-05:00",
  "updated_at": "2020-05-30T16:17:34.8817721-05:00"
}
```
We now have a slice assigned to a subscription. Next we will add a "grain" to that slice.

## Add Grain From File

We will use the `sandpiper` CLI utility to add a test ACES file as a file-based grain. Open a second terminal window (keeping the API server running)

```
open a new terminal window (cd to where your "cli" binary is stored)
./sandpiper -u admin -p admin -c cli-primary.yaml add --slice slice2 testdata/aces-file.xml
```
If you don't see any error, it was added successfully. Note: you can also set environment variables for these user and password parameters.

To see what was added, go back to Insomnia and Send the `List Grains (w/o payload)` request. You should see something like this:
```json
{
  "grains": [
    {
      "id": "793b8629-a237-4e2c-9b35-9d0a8cc4723f",
      "slice_id": "2bea8308-1840-4802-ad38-72b53e31594c",
      "grain_key": "level-1",
      "source": "aces-file.xml",
      "encoding": "z64",
      "payload_len": 2880,
      "created_at": "2020-05-30T21:24:03.97232Z"
    }
  ]
}
```

You can also run the "w/ payload" version to see the ACES file encoded in the "payload". Note that the "encoding" is shown as "z64" which means the file was first zipped and then converted to base64 format for storage in the database.

## Start Secondary Server
In that second terminal window, start a "secondary" server with the following command.
```
open a new terminal window (cd to where your "api" binary is stored)
./api -config api-secondary.yaml

...
Database: "tidepool"
DB Version: 1.15
Server role: "secondary"
Server ID: 20000000-0000-0000-0000-000000000000

⇨ http server started on [::]:8081
```

You should now have two servers listening for commands on separate "ports" and accessing their own data pools. 

Change Insomnia's "Active Environment" from "Primary" (green) to **"Secondary"** (red) using the drop-down menu so we can run API requests against the secondary server. Use the Login request to get an authentication token. Then Add "Company 1".

```
POST Login
POST Add Company 1
```
You should see something like the following:
```json
{
  "id": "10000000-0000-0000-0000-000000000000",
  "name": "Best Brakes",
  "sync_addr": "http://localhost:8080",
  "active": true,
  "created_at": "2020-05-30T16:37:20.7613935-05:00",
  "updated_at": "2020-05-30T16:37:20.7613935-05:00"
}
```

Note this added our trading partner (Best Brakes) with their unique company_id. So we now have the same company in both databases, and their "sync_addr" is pointing to our primary server address (port 8080). In a production environment, this would be a publically accessible url (e.g. sync.bestbrakes.com) where the server is listening.

## Assign API Key

Before we can perform the sync, though, we need to get an API access key from Best Brakes for our sync process. We can get that using the following command (against the primary server). Notice that we changed the user (and password) to "companyadmin".

```
(Activate the Primary (green) environment in Insomnia)
POST Login (using { "username": "companyadmin", "password":"companyadmin" })
POST apikey (from the Sync request folder)
```

It should return something like the following:
```
{
  "primary_id": "10000000-0000-0000-0000-000000000000",
  "sync_api_key": "f1c77a6ee9442d006494d4904476d8f9a328465583cd8e6f3199be99dd5919f41341fc0442b09d297aa33cad..."
}
```
We need to add that key to our secondary database for "Best Brakes". That way, when we initiate the sync, we can pass along the key. We will do that now using Insomnia:
```
(Activate the Secondary (red) environment in Insomnia)
POST Login (using { "username": "admin", "password": "admin" })
PATCH Update apikey (Company 1) (from the Sync request folder)
```
It uses an Insomnia variable we created to copy the generated apikey (from the POST apikey response) to the PATCH request.
 
Normally, this apikey hand-off would be part of the "sign-up" process when a trading partner requests access to your data.

## Sync Secondary with Primary

Now we should be ready to perform the sync. Let's start both a primary and secondary server running (in two separate terminal windows):
```
open a new terminal window (cd to where your "api" binary is stored)
./api -config api-primary.yaml

open a new terminal window (cd to where your "api" binary is stored)
./api -config api-secondary.yaml
```

Open a third terminal window and execute the sync command with the --list option (from the secondary server) to see what trading partners we have defined.

```
./sandpiper -u admin -p admin -c cli-secondary.yaml sync --list

SERVER LIST:
Best Brakes: (http://localhost:8080)
```
Just as we would expect.

Finally, perform the sync process:

```
open a new terminal window (cd to where your "sandpiper" binary is stored)
./sandpiper -u admin -p admin -c cli-secondary.yaml sync
```

If everything worked according to plan, you should see the following results:

```
syncing Best Brakes...
Successful: 1, Errors: 0
```

You can see what was done in Insomnia. Under the Secondary (red) environment, run the List Subscriptions request. These are new subscriptions that were defined on the primary server for eCatCompany. It shows both the company you synced with and the slices that were created.

```json
{
  "subs": [
    {
      "id": "2276e31d-7014-4947-b4ec-1bc170627593",
      "slice_id": "2bea8308-1840-4802-ad38-72b53e31594c",
      "company_id": "10000000-0000-0000-0000-000000000000",
      "name": "Subscription 1 (company2, slice2)",
      "description": "",
      "active": true,
      "created_at": "2020-05-30T22:03:12.039147Z",
      "updated_at": "2020-05-30T22:03:12.039147Z",
      "company": {
        "id": "10000000-0000-0000-0000-000000000000",
        "name": "Best Brakes",
        "sync_addr": "http://localhost:8080",
        "sync_api_key": "76e3ce4cc95b38a1c67ebf398f4c254ee8f173...",
        "active": true,
        "created_at": "2020-05-30T21:37:20.761394Z",
        "updated_at": "2020-05-30T21:58:58.339608Z"
      },
      "slice": {
        "id": "2bea8308-1840-4802-ad38-72b53e31594c",
        "slice_name": "Slice2",
        "slice_type": "aces-file",
        "content_hash": "92891e0296ed93bc03c31e1c09f6287d299431a3",
        "content_count": 1,
        "content_date": "2020-05-30T21:24:03.98232Z",
        "allow_sync": true,
        "sync_status": "success",
        "last_sync_attempt": "2020-05-30T22:03:12.050148Z",
        "last_good_sync": "2020-05-30T22:03:12.07415Z",
        "created_at": "2020-05-30T22:03:12.035152Z",
        "updated_at": "2020-05-30T22:03:12.07415Z"
      }
    }
  ]
}
```

If you use one of the "List Grains" requests, you will also see that the grain was delivered as well.

You've now completed your first sync request.
