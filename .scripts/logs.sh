#!/usr/bin/env bash

docker exec -it restapi tail -f /var/log/supervisor/go-rest-api-err.log
