#!/usr/bin/env bash

SCRIPT=$(dirname "$0")
MIGRATIONS=file://${SCRIPT}/../configs/migrations
CONNECTION=$(yq .db.master "${SCRIPT}/../configs/go-rest-api.conf.yaml")

eval migrate -source "${MIGRATIONS}" -database="${CONNECTION}" up
