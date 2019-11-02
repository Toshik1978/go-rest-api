# Golang RESTful API Example Project
[![Build Status](https://travis-ci.org/Toshik1978/go-rest-api.svg?branch=master)](https://travis-ci.org/Toshik1978/go-rest-api)
[![Coverage](https://codecov.io/gh/Toshik1978/go-rest-api/branch/master/graph/badge.svg)](https://codecov.io/gh/Toshik1978/go-rest-api)
[![License](https://img.shields.io/github/license/toshik1978/go-rest-api)]()
[![Last Commit](https://img.shields.io/github/last-commit/toshik1978/go-rest-api)]()

Project contains some simple Golang application, implemented RESTful service for some abstract financial company.

- [Golang RESTful API Example Project](#golang-restful-api-example-project)
    - [Requirements](#requirements)
        - [Golang](#golang)
        - [Docker](#docker)
        - [Database](#database)
    - [Development](#development)
    - [Quick Start](#quick-start)
        - [In Docker](#in-docker)
        - [Manual](#manual)
    - [Contributing](#contributing)
    - [Additional Documentation](#additional-documentation)
    - [Contacts](#contacts)

## Requirements

### Golang

Project was tested on Go 1.13.3 and 1.13.4, but theoretically should compile on 1.11/1.12 too (all versions with modules support).

### Docker

Docker is not directly required to run application. But for [quick start](#quick-start) you can use it.
So this way you should have `docker` and `docker-compose` [installed on the machine](https://docs.docker.com/install/).

### Database

Project was tested on PostgreSQL 12.0, but theoretically should works fine with old versions too.
You can use docker to run test database instance or you can use your own database installation.

## Development

To develop this project, you should clone it:

```sh
git clone https://github.com/Toshik1978/go-rest-api.git
```

Then you can open it in your preferred editor and have a fun. Project provides `Makefile` to lint/test/build application.
Default `make` command run linter over the project's code, unit tests and finally build binary.
You can manually run linter or linter with tests with the following commands:

```sh
make lint
make test
```

Pay attention, that `test` rule always run `lint` rule.

If you want to skip all stages, except build you can build-only project:

```sh
make build
```

## Quick Start

### In Docker

`.scripts` folder contains some useful shell scripts to quick look into the project.
Default `Dockerfile` and `docker-compose.yml` provides the fastest way to run application.

To instantiate default database instance, run all migrations and start application, you can do:

```sh
.scripts/run.sh
```

Docker allows to run multiple instances of the project. You can add to `docker-compose.yml` additional containers like
`restapi` and map different port to 8080 port inside of container.

### Manual

You can start working with project manually.

First of all you should install dependencies to run migrations:

```sh
make prereq
```

Then install `yq` ([command-line YAML processor](https://yq.readthedocs.io/en/latest/)).
Please, refer documentation for installation details.

With prerequisites installed you can instantiate database. Run:

```sh
.scripts/run.services.sh
```

Then you should run migrations to create database scheme.

```sh
.scripts/migrate.sh
```

Of course you can install database other preferred for you way.
This way you should manually create empty database and update
[application's configuration file](configs/go-rest-api.conf.yaml) to your actual connection string.

Finally you should build project and run it.

```sh
make
./go-rest-api
```

Project contains kind of production configuration file for Docker and `docker-compose-production.yml`.
You can use it instead of previous step with manual build outside of Docker. This way you should build image:

```sh
make docker
```

And then run `docker-compose` with production configuration file.

## Contributing

If you want to contribute into project, you can fork it, change it and create MR. I'll surely look at this!

## Additional Documentation

You can find more information about project in the following documents:

1. [Architecture](docs/architecture.md)
1. [API description](docs/api.md)

## Contacts

You can address any questions on e-mail in public profile.
