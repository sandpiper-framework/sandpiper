# Sandpiper

The Sandpiper framework provides a standard decentralized model to distribute and synchronize shared product
data sets between a primary sender (the "publisher") and a secondary receiver (the "subscriber").

## Getting Started

There are only a few prerequisites for getting a Sandpiper installation up and running. 

#### PostgreSQL Database

Sandpiper requires access to a PostgreSQL sever for its data storage (for both primary and secondary roles). This can be an existing installation (either on premises or in the cloud) or follow
the instructions below to install locally for your desired platform.

Download the binary from the official download site (or use package manager for your platform such as apt or scoop).

[https://www.postgresql.org/download/](https://www.postgresql.org/download/)


See the setup documents in the documentation for platform-specific instructions. Other plug-and-play options are available for those that would like to use pre-configured solutions. Please see the section on Containers and PAAS below. 

#### Sandpiper Binaries

Download the latest Sandpiper Release which contains compiled binaries for both Windows and Linux.

[https://github.com/sandpiper-framework/sandpiper/releases](https://github.com/sandpiper-framework/sandpiper/releases)

There are two programs included in the release, the `api` server and the `sandpiper` command line interface. Put both of them in a directory.

#### Config Files

Both sandpiper programs require configuration settings to run. These settings are stored in [yaml](https://en.wikipedia.org/wiki/YAML) files made up of key/value pairs organized by sections (e.g. Database, Server, Application, etc.). In some cases (such as login credentials), these settings can be overridden by environment variables. See the Deployment section for more information.

Two sample config files are provided (`api-config-sample.yaml` and `cli-config-sample.yaml`) as a template, but these files are also created by the database initialization procedure explained below.

#### Create Database (for each desired server role)

Before we can do anything, we need to create a sandpiper database within the PostgreSQL server. A simple command line tool is provided to take care of this for you. Open a command prompt (terminal) and enter the following commands (assuming you are currently in the root sandpiper folder).

```
cd cmd/cli
./sandpiper init
```
You will be prompted for your PostgreSQL Host Address, Port and Superuser credentials. This is required to create a new database. In most cases, you can simply press Enter to accept the default value (shown in parentheses).

```
PS C:\Users\dougw\autocare\sandpiper\cmd\cli> ./sandpiper init
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
The address `localhost` (which is equivalent to 127.0.0.1) indicates you're running this command on the same machine as PostgreSQL. Otherwise, it would be a standard ip4 address on your network (e.g. 192.168.1.100) or possibly a hosted instance endpoint (e.g. myinstance.123456789012.us-east-1.rds.amazonaws.com). The superuser password (from above) will be hidden when you type.

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

The recommended database name is `sandpiper` regardless of the server-role you require (primary or secondary). If you need both server-roles on one PostgreSQL server, you can name it anything you like ("secondary", "receiver", "tidepool", etc.).

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

## Deployment

This section provides information on how to deploy Sandpiper in a production environment. The system was designed specifically to reduce dependencies to simplify the process. 

### Running in Production

To run Sandpiper in a production environment, simply download the correct api binary from the Sandpiper web site and follow the deployment instructions:

```
./api [-config="path/to/config.yaml"] (defaults to ./config.yaml")
Also supports DB_USER and DB_PASSWORD environment variables
```

### TLS (SSL) Certificate

Discuss how to enable ssl.

### Auto Start -- Linux (the `systemd` init system)

A sample `systemd` "unit" file (sandpiper.service) is included with the executable. Make any required changes and copy to `/etc/systemd/system/`.

Once you enable the service with the following command. It will start automatically on boot, after that.
```
$ sudo systemctl enable sandpiper.service
```
Check status/start/stop/restart
```
$ sudo systemctl {status|start|stop|restart} sandpiper
```
Display all services
```
$ service --status-all
```
https://www.digitalocean.com/community/tutorials/understanding-systemd-units-and-unit-files

### Containers

We provide two pre-configured solutions for testing and deploying a sandpiper installation.

### Docker

https://aranair.github.io/posts/2016/04/27/golang-docker-postgres-digital-ocean/
https://marketplace.digitalocean.com/apps/dokku

### Vagrant

https://www.vagrantup.com/

### PAAS Options

If you would rather not host a sandpiper on your network, there are several ["Platform As A Service"](https://en.wikipedia.org/wiki/Platform_as_a_service) solutions that we support.

#### Render

https://render.com/
https://render.com/docs/deploy-beego
https://render.com/docs/databases

#### Digital Ocean

https://www.digitalocean.com/community/questions/how-to-deploy-golang-program-in-production
https://kenyaappexperts.com/blog/how-to-deploy-golang-to-production-step-by-step/

#### Google App Engine

https://cloud.google.com/appengine/docs/go/

## Source Code

These instructions will help you get up and running on your local machine for development and testing purposes. See the Deployment section for notes on how to deploy the project on a live system.

The following software must be installed on your target development machine.

* [git](https://git-scm.com/downloads) 
* [Go (v13+)](https://golang.org/)
* [PostgreSQL](https://www.postgresql.org/)
* [Task (v2.7+)](https://taskfile.dev/)

### Installing

A step by step series of examples that tell you how to get a development env running

1. Clone the project under the current directory (e.g. $HOME/source/)

```
git clone https://github.com/sandpiper-framework/sandpiper.git
```

## Authentication Endpoints

The application runs as an HTTP server at port 8080. It provides the following RESTful endpoints for authentication:

* `POST /login`: accepts username/passwords and returns jwt token and refresh token
* `GET /refresh/:token`: refreshes sessions and returns jwt token
* `GET /me`: returns info about currently logged in user

An administrator is created as part of the database initialization process. To login to the API, send a POST request to localhost:8080/login with username "admin" and password "admin" in JSON body. **This password must be changed before moving to production.**

Upon a successful login, the response body will include a java web token for subsequent API authentication. These tokens will expire and so must be refreshed (by the client) using the `/refresh` endpoint. 

## Project Structure

1. Root directory contains things not related to code directly, e.g. readme, license, docker-compose, taskfile, etc.

2. Cmd package contains code for starting applications ("main" packages). The directory name for each application should match the name of the executable you want to have. Sandpiper is structured as a monolith application but was written with microservices in mind. We use the Go convention of placing each main package as a subdirectory of the cmd directory. As an example, the "client" application's binary would be located under `cmd/client`. It also loads the necessary configuration and passes it to the service initializers.

3. The rest of the code is located under two directories: `shared` and `pkg`. The pkg directory contains a directory for each of the executables found in the cmd directory. These can be thought of like microservices of Sandpiper. The shared directory contains code common to all cmds. 

4. Microservice directories, like api (naming corresponds to `cmd/` folder naming) contains multiple folders for each domain it interacts with, for example: `user`, `company`, `slice` etc.

5. Domain directories, like "user", contain all application/business logic and three additional directories: "logging", "platform" and "transport".

6. The `platform` folder contains various packages that provide support for things like databases, authentication and marshaling. Most of the packages located under platform are decoupled by using interfaces (with the "Repository pattern"). Every platform has its own package, for example, pgsql (orm for postgres), elastic, redis, memcache etc.

7. The `transport` package contains HTTP handlers. This package receives the requests, marshals, validates then passes it to the corresponding service.

8. The `internal` folder contains helper packages and models. Packages such as mock, middleware, configuration, server are located here.

## Running the tests

... Explain how to run automated tests ...

## Development

This section highlights areas helpful for continued development of the project.

### Calling API Endpoints Using Insomnia

An [Insomnia](https://insomnia.rest/) workspace configuration is provided in the `/test/api` directory with instructions in the
associated README.md file.

### Implementing CRUD of another table

Let's say you have a table named 'cars' that handles employee's cars. To implement CRUD on this table you need:

1. Inside `pkg/shared/model` create a new file named `car.go`. Inside put your entity (struct), and methods on the struct if you need them.

2. Create a new `car` folder in the (micro)service where your service will be located, most probably inside `api`. Inside create a file/service named car.go and test file for it (`car/car.go` and `car/car_test.go`). You can test your code without writing a single query by mocking the database logic inside /mock/mockdb folder. If you have complex queries interfering with other entities, you can create in this folder other files such as car_users.go or car_templates.go for example.

3. Inside car folder, create folders named `platform`, `transport` and `logging`.

4. Code for interacting with a platform like database (postgresql) should be placed under `car/platform/pgsql`. (`pkg/api/car/platform/pgsql/car.go`)

5. In `pkg/api/car/transport` create a new file named `http.go`. This is where your handlers are located. Under the same location create http_test.go to test your API.

6. In logging directory create a file named `car.go` and copy the logic from another service. This serves as request/response logging.

7. In `pkg/api/api.go` wire up all the logic, by instantiating car service, passing it to the logging and transport service afterwards.

### Implementing other platforms

Similar to implementing APIs relying only on a database, you can implement other platforms by:

1. In the service package, in car.go add interface that corresponds to the platform, for example, Indexer or Reporter.

2. Rest of the procedure is same, except that in `/platform` you would create a new folder for your platform, for example, `elastic`.

3. Once the new platform logic is implemented, create an instance of it in main.go (for example `elastic.Client`) and pass it as an argument to car service (`pkg/api/car/car.go`).

### Running database queries in a transaction

To use a transaction, before interacting with db create a new transaction:

```
err := s.db.RunInTransaction(func (tx *pg.Tx) error{
    // Application service here
})
```

Instead of passing database client as `s.db` , inside this function pass it as `tx`. Handle the error accordingly.

## Contributing

Please read [CONTRIBUTING](CONTRIBUTING.md) for details on our process for submitting pull requests. Note also that this project is released with a Contributor Code of Conduct. By participating in this project you agree to abide by its terms.

[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-v2.0%20adopted-ff69b4.svg)](code_of_conduct.md)


## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/your/project/tags). 

## Authors

* **Doug Winsby** - *Initial work* - [Winsby Group LLC](https://winsbygroup.com)

See also the list of [contributors](https://github.com/orgs/sandpiper-framework/people) who participated in this project.

## License & Copyrights

Copyright (c) 2019-2020 The Sandpiper Authors. All rights reserved.

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

The Sandpiper logo and mascot images are Copyright 2020 [Megan Winsby](https://www.linkedin.com/in/mwinsby/). Used with permission.

## Acknowledgments

1. [Echo](https://echo.labstack.com/) - HTTP 'framework'.
2. [Go-Pg](https://github.com/go-pg/pg) - PostgreSQL ORM
3. [JWT-Go](https://github.com/dgrijalva/jwt-go) - JWT Authentication
4. [Zerolog](https://github.com/rs/zerolog) - Structured logging
5. [Bcrypt](https://github.com/golang/crypto/) - Password hashing
6. [Yaml](https://gopkg.in/yaml.v2) - Unmarshalling YAML config file
7. [Validator](https://github.com/go-playground/validator) - Request validation.
8. [lib/pq](https://github.com/lib/pq) - PostgreSQL driver
9. [zxcvbn-go](https://github.com/nbutton23/zxcvbn-go) - Password strength checker
10. [DockerTest](https://github.com/fortytw2/dockertest) - Testing database queries (might need to change this lib choice)
11. [Testify/Assert](https://github.com/stretchr/testify) - Asserting test results
12. [go-bindata](https://github.com/go-bindata/go-bindata) - Turn data file into go code (for migrations)
13. [uuid](https://github.com/google/uuid) - Google's library to generate and manipulate uuid values
14. [cli](https://github.com/urfave/cli) - Command line argument processing

Efforts were made to abstract and localize these dependencies.

### Notices required by included works

**Starter template** ([GORSK](https://github.com/ribice/gorsk)) provided by:

Copyright (c) 2018 Emir Ribić

<small>Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.
