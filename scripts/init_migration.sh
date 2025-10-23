#!/bin/bash

NAME=$1
PWD=$(pwd)

if [ -z "$NAME" ]; then
  echo "Usage: ./init_migration.sh migration_name"
  exit 1
fi

docker run --rm \
  -v "$PWD/../backend/migrations:/migrations" \
  migrate/migrate:latest \
  create -ext sql -dir /migrations -seq "$NAME"
