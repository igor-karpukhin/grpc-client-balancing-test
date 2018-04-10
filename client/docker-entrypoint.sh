#!/bin/sh
set -e

: ${ADDR:=127.0.0.1:9090}

if [ "$1" = 'test-gclient' ]; then
    exec /usr/local/bin/test-gclient -addr=${ADDR}
fi

exec "$@"
