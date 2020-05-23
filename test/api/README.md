# Testing Handbook

This handbook walks you through the testing process. Note that all directory paths are shown in Linux format (forward slash) since they also work in Windows PowerShell. Also, all paths (and so `cd` commands) are shown relative to the root `sandpiper` directory.

## Download Sandpiper Distribution Package

https://github.com/sandpiper-framework/sandpiper/releases

## Download and Install PostgreSQL

https://www.postgresql.org/download/

See our separate guides for installing PostgreSQL on specific platforms.

Be sure to write down the superuser (usually `postgres`) and password. You will need them to manage your PostgreSQL server.

## Creating the Primary Database

Before we can do anything, we need to create a sandpiper database within the PostgreSQL server. A simple command line tool is provided to take care of this for you. Open a command prompt (terminal) window and enter the following commands (assuming you are currently in the root sandpiper folder).

```
cd cmd/cli
./sandpiper init
```
You will be prompted for your PostgreSQL Host Address, Port and Superuser credentials. This is required to create a new database. In most cases, you can simply press Enter to accept the default value (shown in parentheses).

```
PS C:\Users\dougw\autocare\sandpiper\cmd\cli> ./sandpiper init --id 10000000-0000-0000-0000-000000000000
sandpiper (v0.1.2-67-g5facfce-dirty)
Copyright 2020 The Sandpiper Authors. All rights reserved.

INITIALIZE A SANDPIPER DATABASE

PostgreSQL Address (localhost):
PostgreSQL Port (5432):
PostgreSQL Superuser (postgres):
PostgreSQL Superuser Password: *********
SSL Mode (disable):
connected to host
```
Notice that we included the `--id` option on the init command. This option lets you provide the `server_id` rather than having the software assign a random unique value (allowing our existing tests to work without change). `localhost` (which is equivalent to 127.0.0.1) indicates you're running this command on the same machine as PostgreSQL. Otherwise, it would be a standard ip4 address on your network (e.g. 192.168.1.100) or possibly a hosted instance endpoint (e.g. myinstance.123456789012.us-east-1.rds.amazonaws.com). The superuser password (from above) will be hidden when you type.

You should see "connected to host" to indicate that the connection was successful. Next, you will be prompted for the new database information.

```
New Database Name (sandpiper):
Database Owner (sandpiper):
Database Owner Password: focal-weedy-brood-hat
CREATE DATABASE sandpiper;
CREATE USER sandpiper WITH ENCRYPTED PASSWORD 'focal-weedy-brood-hat';
user "sandpiper" already exists
GRANT ALL PRIVILEGES ON DATABASE sandpiper TO sandpiper;

applying migrations...
Database: "sandpiper"
DB Version: 1 (migrated from 0 to 1)
```

The recommended database name is `sandpiper` regardless of the server-role you require (primary or secondary). If you need both server-roles, it's customary to call the secondary database `tidepool` (but any name will suffice).

In the example above, we used default values except when required to enter a password for the database owner (please select a strong password of your own) and again, keep a record of it for later. You will need it when starting the server.

The database owner is the only user to connect directly to the database (via the sandpiper-api server). This should not be confused with a sandpiper end-user which is stored in the `users` table for authentication and access.

```
Company Name: Better Brakes
Server-Role (primary*/secondary): primary
Public Sync URL: http://localhost:8081
Added Company "Better Brakes"
Sandpiper Admin Password: admin
Added User "admin"

initialization complete for "sandpiper"
```

In production, you would enter a strong admin password, but enter "admin" here to make testing easier. Also, the public sync URL would normally be something like `https://sandpiper.betterbrakes.com`, but we are going to test locally, on the same machine with both servers.
<div class="page"/>

## Creating the Secondary Database

Use the same procedure as above, but use this command to set the receiver's server-id:
 
 ```
 ./sandpiper init --id 20000000-0000-0000-0000-000000000000
 ```
 
Enter `tidepool` for the database name, and because it is a secondary (i.e. receiver), it will not prompt for a public sync URL.

```
Company Name: eCat Co
Server-Role (primary*/secondary): secondary
Added Company "eCat Co"
Sandpiper Admin Password: admin
Added User "admin"

initialization complete for "tidepool"
```

You should now have two databases, each with an "admin" user and an associated "company".

Next we'll run the sandpiper server (using the "primary" database) and create `subscriptions` and `grains` for us to sync. We'll do most of this work with a free REST client called Insomnia (someone must have thought that name was clever).

## Insomnia REST Client

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
cd cmd/api
copy config-sample.yaml primary.yaml
notepad primary.yaml
(make any necessary changes to the `database` section and save the file) 
./api --config primary.yaml
```
If everything is working properly, you should see something like the following:

```
Sandpiper API Server (v0.1.2-67-g5facfce-dirty)
Copyright 2019-2020 The Sandpiper Authors. All rights reserved.

Database: "sandpiper"
DB Version: 1
Server role: "primary"

⇨ http server started on [::]:8080
```
That last line shows that a web server is running and listening on http://localhost:8080. You should be able to open a browser and type that address and receive a response of:

```
"Sadnpiper API OK"
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
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjIjoiNWMzN2QyZmQtY2EzZC00ZTVlLThkOGEtNDUwNmFlMzA1OGUwIiwiZSI6ImFkbWlu...",
  "expires": "2020-05-18T16:27:43-05:00",
  "refresh_token": "d16e1a83587fd1546a204725c97dcca29c2c7422"
}
```

That token has now been saved in our `jwt` variable, and we can use it to make authenticated requests (as admin). The user information is actually encrypted in the token itself so that the server knows who is making those requests.

If you look at the console where you're running the server, you will also notice that there is lots of information displayed whenever a request is made. There's a setting in the config file to supress some of these messages, but for now it can be helpful while testing to keep it.

## Companies and Users

Let's move on to adding some Companies and Users. Select the `List Companies` request and press `Send`. You should see something like the following:

```json
{
  "companies": [
    {
      "id": "10000000-0000-0000-0000-000000000000",
      "name": "Better Brakes",
      "sync_addr": "http://localhost:8081",
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
  "page": 0
}
```

This is the primary server's owner. It currently has just one user (which is also the one you logged in under).

It's interesting to note that the company `id` is a globally unique identifier (because they can be shared across servers), but the user `id` is simply a sequential number (because they are local to this server).

A good example of a shared company is "eCat Co" (which we created as the owner of our secondary server above). When a primary company agrees to share its product data, it adds a "secondary" company to its list of companies it does business with.
 
Let's do that now by selecting the `Add Copmany (eCat)` request and pressing `Send`. Then, let's add a user for our eCat company (`Add User (eCat)` under the User Resource). If you now go back and send the List Companies request again, you should see this second one added to the list:

```json
{
  "id": "20000000-0000-0000-0000-000000000000",
  "name": "eCat Co",
  "sync_addr": "http://sync.ecatco.com",
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

Now that we have two trading partners defined, we need a structure that organizize the data we want to share between them. This structure is called a `slice` and the way we assign these slices to companies is with a `subscription`.

Under the "Slice Resource" folder you should see two "POST Add" requests. In each case, the request body is provided to create a new slice.
