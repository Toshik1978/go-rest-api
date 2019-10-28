#!/usr/bin/env bash

SCRIPT=$(dirname "$0")
MIGRATIONS=${SCRIPT}/../configs/migrations

if [ -z "$1" ]; then
  echo "Migration name should be passed"
  exit 255
fi

migrate create -ext sql -dir "${MIGRATIONS}" -seq $1
