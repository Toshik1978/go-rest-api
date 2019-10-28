# Golang RESTful API example project

Project contains some simple Golang application, implemented RESTful service for some abstract financial company.

## Requirements

### Golang

Project was tested on Go 1.13.3, but theoretically should compile on 1.11/1.12 too (all versions with modules support).

### Docker

Docker is not directly required to run application. But for quick start you can use `docker-compose`.
So this way you should have `docker` and `docker-compose` [installed on the machine](https://docs.docker.com/install/).

### Database

#### PostgreSQL

Project was tested on PostgreSQL 12.0, but theoretically should works fine with old versions too.
You can use `docker-compose` to run test database instance. Or you can use your own database installation.
This way you should update [application's configuration file](configs/go-rest-api.conf.yaml) to your actual connection string.

#### Migrations

To run migration you should install `yq` ([command-line YAML processor](https://github.com/mikefarah/yq)).
After tool installation, PostgreSQL running and updating configuration file to actual connection string,
you can run migration scripts to update database.

```sh
.scripts/migrate.sh
```

## Running

You can run sample instance via `docker-compose` with [provided](docker-compose.yml) `docker-compose.yml`.
