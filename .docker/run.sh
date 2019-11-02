#!/usr/bin/env sh

# Run migrations
SCRIPT=$(dirname "$0")
MIGRATIONS=${SCRIPT}/migrations
MIGRATIONS_SCHEME=file://${MIGRATIONS}
CONNECTION=$(yq .db.master "${SCRIPT}/configs/go-rest-api.conf.yaml")

MIGRATIONS_COUNT=$(ls -1 ${MIGRATIONS} | wc -l)
if [ "${MIGRATIONS_COUNT}" -gt 0 ]; then
  eval /opt/go-rest-api/migrate -source "${MIGRATIONS_SCHEME}" -database="${CONNECTION}" up
fi

# Run daemon
/usr/bin/supervisord -c /etc/supervisord.conf
