#!/bin/sh

PROCESS=${FLY_PROCESS_GROUP}

if [ "$PROCESS" = "server" ]; then
  echo "Starting server process"
  exec litefs mount -- $@
elif [ "$PROCESS" = "worker" ]; then
  echo "Starting worker process"
  exec $@
else
  echo "Unknown process"
  exit 1
fi
