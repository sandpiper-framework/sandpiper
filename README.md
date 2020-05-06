# Sandpiper

The Sandpiper framework provides a standard decentralized model to classify, distribute, and synchronize shared product
data sets between an originating sender (the "publisher") and a secondary receiver (the "subscriber").

## Getting Started

There are only a few prerequisites for getting a Sandpiper installation up and running. 

#### PostgreSQL Database

Sandpiper requires access to a PostgreSQL sever for its data storage (for both primary and secondary roles). This can be an existing installation (either on premises or in the cloud) or follow
the instructions below to install locally for your desired platform.

Download binary from official download site (or use package manager for your platform such as apt or scoop).

```
https://www.postgresql.org/download/
```

##### Windows

To start/stop service, run `pg_ctl start`, `pg_ctl stop`.

##### Linux

To start the service, run

```
sudo service postgresql start
```

##### Create Database (for each desired role)

*todo:* Database creation is automated with 'task', but we don't want that dependency for standard production installations. Create .bat and .sh files to handle this for those not running from source.

```
task create-pub-db | create-sub-db | create-both-db
```

```
win64: psql --username=postgres --file=db_create.sql
linux: sudo -u postgres psql --username=postgres --file=db_create.sql
```

Enter the master postgresql user password when prompted. The following commands should be included in the db_create.sql file (depending on desired role).

*Primary:*
```
CREATE DATABASE sandpiper;
CREATE USER admin WITH ENCRYPTED PASSWORD '--password here--';
GRANT ALL PRIVILEGES ON DATABASE sandpiper TO admin;
```

*Secondary:*
```
CREATE DATABASE tidepool;
CREATE USER admin WITH ENCRYPTED PASSWORD '--password here--';
GRANT ALL PRIVILEGES ON DATABASE tidepool TO admin;
```

### Running in Production

To run Sandpiper in a production environment, simply download the correct binary from the Sandpiper web site and follow the deployment instructions:

[Downloads](https://sandpiper.org/downloads)

```
./api [-config="path/to/config.yaml"] (defaults to ./config.yaml")
Also supports DB_USER and DB_PASSWORD environment variables
```
## Deployment

Add additional notes about how to deploy this on a live system

### Linux (the `systemd` init system)

A sample systemd "unit" file (sandpiper.service) is included with the executable. Make any required changes and copy to `/etc/systemd/system/`.

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

### Source Code

These instructions will help you get up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

The following software must be installed on your target development machine.

* [git](https://git-scm.com/downloads) 
* [Go](https://golang.org/)
* [PostgreSQL](https://www.postgresql.org/)
* [Task](https://taskfile.dev/)

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

Please read [CONTRIBUTING](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/your/project/tags). 

## Authors

* **Doug Winsby** - *Initial work* - [Winsby Group LLC](https://winsbygroup.com)

See also the list of [contributors](https://github.com/orgs/sandpiper-framework/people) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

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
