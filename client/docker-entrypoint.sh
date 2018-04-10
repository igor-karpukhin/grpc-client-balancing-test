#!/bin/sh
set -e

: ${ADDR:=127.0.0.1:9090}
: ${DNS:=false}

if [ "$1" = 'test-gclient' ]; then
    exec /usr/local/bin/test-gclient -addr=${ADDR} -dns=${DNS}
fi

exec "$@"
