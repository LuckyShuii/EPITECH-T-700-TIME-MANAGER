#!/bin/sh

# To use this script, you need to pass one of the following commands as an argument:
# up: to apply all up migrations
# down: to apply all down migrations
# version: to print the current migration version
# force VERSION: to set the migration version without running migrations

# Example usage:
# ./migrate.sh up
# ./migrate.sh down
# ./migrate.sh version
# ./migrate.sh force 5

# Generate a migration in CLI:
# docker run --rm -v backend/migrations:/migrations migrate/migrate create -ext sql -dir /migrations -seq create_users_table

COMMAND=$1
ARG=$2

if [ "$COMMAND" = "force" ] && [ -n "$ARG" ]; then
  COMMAND="force $ARG"
fi

if [ -z "$COMMAND" ]; then
  echo "Usage: ./migrate.sh [up|down|version|force VERSION]"
  exit 1
fi

. ../.env

NETWORK_NAME="t-dev-700-project-par_19_default"

echo "📡 Using network: $NETWORK_NAME"
echo "🔗 postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@database:5432/${POSTGRES_DB}?sslmode=disable ($COMMAND)"

# If the command is "down", we want to run "down -all" to drop all tables
if [ "$COMMAND" = "down" ]; then
  COMMAND="down -all"
fi

docker run --rm \
  -v $(pwd)/../backend/migrations:/migrations \
  --network $NETWORK_NAME \
  migrate/migrate:latest \
  -source=file:///migrations \
  -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@database:5432/${POSTGRES_DB}?sslmode=disable" \
  $COMMAND
