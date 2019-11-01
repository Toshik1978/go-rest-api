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
    - [Build](#build)
    - [Quick Start](#quick-start)
        - [In Docker](#in-docker)
        - [Manual](#manual)
    - [Contributing](#contributing)
    - [Additional Documentation](#additional-documentation)
    - [Contacts](#contacts)

## Requirements

### Golang

Project was tested on Go 1.13.3, but theoretically should compile on 1.11/1.12 too (all versions with modules support).

### Docker

Docker is not directly required to run application. But for quick start you can use `docker-compose`.
So this way you should have `docker` and `docker-compose` [installed on the machine](https://docs.docker.com/install/).

### Database

Project was tested on PostgreSQL 12.0, but theoretically should works fine with old versions too.
You can use `docker-compose` to run test database instance.

Also you can use your own database installation. This way you should manually create empty database and update
[application's configuration file](configs/go-rest-api.conf.yaml) to your actual connection string.

To run migrations you should install `yq` ([command-line YAML processor](https://yq.readthedocs.io/en/latest/)). Please,
refer documentation for installation details.

## Build

Default `make` command run linter over the project's code, unit tests and finally build binary.
You can manually run linter or linter with tests with the following commands:

```sh
make lint
make test
```

Additional you can build-only project:

```sh
make build
```

## Quick Start

### In Docker

You can run sample instance with help of `docker-compose` with [provided](docker-compose.yml) `docker-compose.yml`. Or

```sh
.scripts/run.sh
```

Then you should run migrations to create database scheme.

```sh
.scripts/migrate.sh
```

Docker allows to run multiple instances of the project. You can add to `docker-compose.yml` additional containers like
`restapi` and map different port to 8080 port inside of container.

### Manual

You can run PostgreSQL instance with help of `docker-compose` with [provided](docker-compose.yml) `docker-compose.yml`. Or 

```sh
.scripts/run.services.sh
```

Then you should run migrations to create database scheme.

```sh
.scripts/migrate.sh
```

Finally you should [build project](#build) and run it.

```sh
./go-rest-api
```

## Contributing

If you have any plans to contribute this project, you can fork it, change it and create MR. I'll surely look at this!

## Additional Documentation

You can find more information about project in the following documents:

1. [Architecture](docs/architecture.md)
1. [API description](docs/api.md)

## Contacts

You can address any questions to my e-mail in public profile.
