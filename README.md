# Sandpiper

The Sandpiper framework provides a standard decentralized model to classify, distribute, and synchronize shared product
data sets between an originating sender (the "publisher") and a derivative receiver (the "subscriber").

## Getting Started

There are only a few prerequisites for getting a Sandpiper installation up and running. 

#### PostgreSQL Database

Sandpiper requires access to a PostgreSQL sever for its data storage (for both primary and secondary roles). This can be an existing installation (either on premises or in the cloud) or follow
the instructions below to install locally for your desired platform.

##### Windows

https://www.postgresql.org/download/windows/

To start/stop service, run `pg_ctl start`, `pg_ctl stop`.

##### Linux

https://www.postgresql.org/download/windows/

To start/stop service, run `pg_ctl start`, `pg_ctl stop`.

##### Create Database (for each desired role)

*todo:* this process is automated with 'task', but we don't want that dependency. Create .bat and .sh files to handle this for those not running from source.

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

To run Sandpiper in a production environment, simply download the correct binary from the Sandpiper web site and follow the installation instructions:

[Downloads](https://sandpiper.org/downloads)

### Source Code

These instructions will help you get up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

The following software must exist on your target development machine.

* [Go](https://golang.org/)
* [PostgreSQL](https://www.postgresql.org/)
* [Task](https://taskfile.dev/)

### Installing

A step by step series of examples that tell you how to get a development env running

```
Give the example
```

End with an example of getting some data out of the system or using it for a little demo

## Endpoints

The application runs as an HTTP server at port 8080. It provides the following RESTful endpoints:

* `POST /login`: accepts username/passwords and returns jwt token and refresh token
* `GET /refresh/:token`: refreshes sessions and returns jwt token
* `GET /me`: returns info about currently logged in user
* `GET /v1/users`: returns list of users
* `GET /v1/users/:id`: returns single user
* `POST /v1/users`: creates a new user
* `PATCH /v1/password/:id`: changes password for a user
* `DELETE /v1/users/:id`: deletes a user

You can log in as admin to the application by sending a post request to localhost:8080/login with username `admin` and password `admin` in JSON body.


## Project Structure

1. Root directory contains things not related to code directly, e.g. docker-compose, CI/CD, readme, bash scripts etc. It may also contain vendor folder.

2. Cmd package contains code for starting applications (main packages). The directory name for each application should match the name of the executable you want to have. Sandpiper is structured as a monolith application but can be restructured to contain multiple microservices. We use the Go convention of placing main package as a subdirectory of the cmd package. As an example, scheduler application's binary would be located under cmd/cron. It also loads the necessary configuration and passes it to the service initializers.

3. Rest of the code is located under /pkg. The pkg directory contains `utl` and 'microservice' directories.

4. Microservice directories, like api (naming corresponds to `cmd/` folder naming) contains multiple folders for each domain it interacts with, for example: user, car, appointment etc.

5. Domain directories, like "user", contain all application/business logic and two additional directories: "platform" and "transport".

6. Platform folder contains various packages that provide support for things like databases, authentication and marshaling. Most of the packages located under platform are decoupled by using interfaces. Every platform has its own package, for example, postgres, elastic, redis, memcache etc.

7. Transport package contains HTTP handlers. The package receives the requests, marshals, validates then passes it to the corresponding service.

8. Utl directory contains helper packages and models. Packages such as mock, middleware, configuration, server are located here.

## Running the tests

Explain how to run the automated tests for this system

### Break down into end to end tests

Explain what these tests test and why

```
Give an example
```

## Deployment

Add additional notes about how to deploy this on a live system

## Built With

* [Dropwizard](http://www.dropwizard.io/1.0.2/docs/) - The web framework used
* [Maven](https://maven.apache.org/) - Dependency Management
* [ROME](https://rometools.github.io/rome/) - Used to generate RSS Feeds

## Contributing

Please read [CONTRIBUTING.md](https://gist.github.com/PurpleBooth/b24679402957c63ec426) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/your/project/tags). 

## Authors

* **Billie Thompson** - *Initial work* - [PurpleBooth](https://github.com/PurpleBooth)

See also the list of [contributors](https://github.com/your/project/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

* Hat tip to anyone whose code was used
* Inspiration
* etc


1. Echo - HTTP 'framework'.
2. Go-Pg - PostgreSQL ORM
3. JWT-Go - JWT Authentication
4. Zerolog - Structured logging
5. Bcrypt - Password hashing
6. Yaml - Unmarshalling YAML config file
7. Validator - Request validation.
8. lib/pq - PostgreSQL driver
9. zxcvbn-go - Password strength checker
10. DockerTest - Testing database queries
11. Testify/Assert - Asserting test results

Most of these can easily be replaced with your own choices since their usage is abstracted and localized.
